#!/usr/bin/env ts-node
/* 
清理目录 ，这个脚本会直接删除某些目录
*/
import shell from 'shelljs';

if (!shell.which('git')) {
  shell.echo('Sorry, this script requires git');
  shell.exit(1);
}

// 清理目录
shell.exec('git clean -fdX');
// 进程 结束
process.exit(0);
