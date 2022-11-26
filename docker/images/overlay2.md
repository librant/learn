
1、 overlay2运作方式
- merged
- lowerdir
- upperdir
- workdir

- lowerdir=<dir>：指定用户需要挂载的lower层目录，lower层支持多个目录，用“:”间隔，优先级依次降低。最多支持500层。
- upperdir=<dir>：指定用户需要挂载的upper层目录，upper层优先级高于所有的lower层目录。
- workdir=<dir>：指定文件系统挂载后用于存放临时和间接文件的工作基础目录。
- default_permissions：
- redirect_dir=on/off：开启或关闭redirect directory特性，开启后可支持merged目录和纯lower层目录的rename/renameat系统调用。
- index=on/off：开启或关闭index特性，开启后可避免hardlink copyup broken问题。