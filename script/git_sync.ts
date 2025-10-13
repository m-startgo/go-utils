import os from 'os';
import { $ } from 'bun';

// 获取当前执行命令的目录
const currentDir = process.cwd();
console.info('当前执行命令的目录:', currentDir);
$.cwd(currentDir);

const SetFileMod777 = async () => {
  const sysType = os.platform();
  if (sysType == 'darwin' || sysType == 'linux') {
    await $`chmod -R 777 ./`;
    console.info('文件权限已重置');
  }
};

const SetGitLocalConfig = async () => {
  try {
    // 开启当前项目 git 大小写敏感
    await $`git config core.ignorecase false`;
    // 忽略当前项目文件权限变更
    await $`git config core.filemode false`;
    // 禁用当前项目 pull.rebase
    await $`git config pull.rebase false`;
    console.info('本地 git config 已覆盖');
  } catch (error) {
    console.error(`git err code: ${error.exitCode}`);
    console.info(error.stdout.toString());
    console.info(error.stderr.toString());
    process.exit(1);
  }
};

await SetFileMod777();
await SetGitLocalConfig();

const desc = process.argv[2];
if (desc) {
  console.info(`git commit: ${desc} \n`);
} else {
  console.warn(`请认真填写 git commit ! \n`);
  process.exit(0);
}

try {
  await $`git pull`;
  await $`git add .`;
  await $`git commit -m "${desc}"`;
  await $`git push`;
} catch (error) {
  if (error.stdout.toString().includes('nothing to commit, working tree clean')) {
    process.exit(0);
  } else {
    console.error(`git err code: ${error.exitCode}`);
    throw new Error(error.stdout.toString());
  }
}
process.exit(0);
