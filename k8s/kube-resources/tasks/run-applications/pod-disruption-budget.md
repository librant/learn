
#### Pod 中断预算 （Pod Disruption Budget）

1) 终端预算的工作原理   
   应用程序所有者可以为每个应用程序创建一个 PodDisruptionBudget 对象（PDB）；
   PDB 将限制在同一时间自愿中断的复制应用程序中宕机的 Pod 的数量；

