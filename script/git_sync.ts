#!/usr/bin/env ts-node

/* 
自动提交脚本，只提交道当前分支的远端仓库
*/

import os from 'os';
import shell from 'shelljs';

if (!shell.which('git')) {
  shell.echo('Sorry, this script requires git');
  shell.exit(1);
}

const desc = process.argv[2];
if (!desc) {
  console.error('err:请提供提交描述信息');
  process.exit(0);
}

// 获取当前执行命令的目录
const currentDir = process.cwd();
shell.cd(currentDir);
console.info('当前执行命令的目录:', shell.pwd().toString());

// 使用 os 设置所有文件权限为 777
function SetFileMod777() {
  if (os.platform() === 'win32') {
    console.info('Windows 平台，跳过权限设置');
    return;
  }
  shell.exec('chmod -R 777 ./');
  console.info('已将所有文件权限设置为 777');
}

const SetGitLocalConfig = async () => {
  shell.exec('git config core.ignorecase false');
  shell.exec('git config core.filemode false');
  shell.exec('git config pull.rebase false');
};

await SetFileMod777();
await SetGitLocalConfig();

shell.exec('git pull');
shell.exec('git add .');
const commitCmd = `git commit -m "${desc}"`;
const commitResult = shell.exec(commitCmd);

if (commitResult.code === 0) {
  shell.exec('git push');
} else if (
  commitResult.code === 1 &&
  commitResult.stdout.includes('nothing to commit, working tree clean')
) {
  console.info('No changes to commit. Exiting.');
  shell.exit(0);
} else {
  console.error(`Git commit failed with code: ${commitResult.code}`);
  console.error(commitResult.stdout);
  shell.exit(1);
}

process.exit(0);
