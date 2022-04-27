## Cloud Native Some Knowledge Summary

云原生技术有利于各组织在公有云、私有云和混合云等新型动态环境中，构建和运行可弹性扩展的应用。云原生的代表技术包括容器、服务网格、微服务、不可变基础设施和声明式API。

这些技术能够构建容错性好、易于管理和便于观察的松耦合系统。结合可靠的自动化手段，云原生技术使工程师能够轻松地对系统作出频繁和可预测的重大变更。

云原生计算基金会（CNCF）致力于培育和维护一个厂商中立的开源生态系统，来推广云原生技术。我们通过将最前沿的模式民主化，让这些创新为大众所用。

Reference link: https://github.com/cncf/toc/blob/main/DEFINITION.md

### containerization

容器技术的核心功能，就是通过约束和修改进程的动态表现，从而为其创造出一个“边界”。对于 Docker 等大多数 Linux 容器来说，Cgroups 技术是用来制造约束的主要手段，而 Namespace 技术则是用来修改进程视图的主要方法。

Docker 容器这个听起来玄而又玄的概念，实际上是在创建容器进程时，指定了这个进程所需要启用的一组 Namespace 参数。
这样，容器就只能“看”到当前 Namespace 所限定的资源、文件、设备、状态，或者配置。而对于宿主机以及其他不相关的程序，它就完全看不到了。
所以说，容器，其实是一种特殊的进程而已。

Linux Cgroups 的全称是 Linux Control Group。它最主要的作用，就是限制一个进程组能够使用的资源上限，包括 CPU、内存、磁盘、网络带宽等等。

Linux Cgroups 的设计还是比较易用的，简单粗暴地理解呢，它就是一个子系统目录加上一组资源限制文件的组合。而对于 Docker 等 Linux 容器项目来说，它们只需要在每个子系统下面，为每个容器创建一个控制组（即创建一个新目录），然后在启动容器进程之后，把这个进程的 PID 填写到对应控制组的 tasks 文件中就可以了。


这个挂载在容器根目录上、用来为容器进程提供隔离后执行环境的文件系统，就是所谓的“容器镜像”。它还有一个更为专业的名字，叫作：rootfs（根文件系统）。


对 Docker 项目来说，它最核心的原理实际上就是为待创建的用户进程：
启用 Linux Namespace 配置；
设置指定的 Cgroups 参数；
切换进程的根目录（Change Root）。

这样，一个完整的容器就诞生了。不过，Docker 项目在最后一步的切换上会优先使用 pivot_root 系统调用，如果系统不支持，才会使用 chroot。这两个系统调用虽然功能类似，但是也有细微的区别，这一部分小知识就交给你课后去探索了。

另外，需要明确的是，rootfs 只是一个操作系统所包含的文件、配置和目录，并不包括操作系统内核。在 Linux 操作系统中，这两部分是分开存放的，操作系统只有在开机启动时才会加载指定版本的内核镜像。
所以说，rootfs 只包括了操作系统的“躯壳”，并没有包括操作系统的“灵魂”。
那么，对于容器来说，这个操作系统的“灵魂”又在哪里呢？
实际上，同一台机器上的所有容器，都共享宿主机操作系统的内核。



Namespace

/proc/[process id]/ns/

PID, UTS, Network, User, Mount, IPC, CGroup 

lsns
lsns -t net
ls -la /proc/<pid>/ns/
nsenter -t <pid> -x xxx    --- > nsenter -t <pid> -n ip a #go to net ns, execute ip a command


CGroup


/sys/fs/cgroup/


AUFS

/var/lib/docker/aufs/diff
/var/lib/docker/aufs/mnt


  
### Modernization
  
  
  
### Kubernetes

架构设计原则

apiserver


控制器协同工作
apiserver和etcd通信，apiserver提供watch接口给到其他组件，所有组件只和apiserver通信。
scheduler根据由kubelet汇报上来的节点信息(在etcd)， 根据算法将pod与某节点绑定 -- 更新pod对象nodeName信息在etcd, kubelet监听到此消息，负责pod生命周期管理。

kubectl client -- > apiserver -- > etcd -- > deployment controller -- > create rs -- > replicaset controller -- > create pods -- > scheduler bind pods with nodes -- > kubelet create pods

kubectl client/deployment controller/rs controller/scheduler/kubelet都是和apiserver打交道

ETCD

CRI/CNI/CSI

kube-proxy

coredns


scheduler


网络


集群联邦




### Service Mesh

envoy

istio流量劫持

istio traffic management

istio mTLS

istio authentication/authorization



### DevOps
Through six years of research, the DevOps Research and Assessment (DORA) team has identified four key metrics that indicate the performance of software delivery. Four Keys allows you to collect data from your development environment (such as GitHub or GitLab) and compiles it into a dashboard displaying these key metrics.

These four key metrics are:

Deployment Frequency
Lead Time for Changes
Time to Restore Services
Change Failure Rate  

Reference link: https://github.com/GoogleCloudPlatform/fourkeys


职责：  
大规模运行和管理您的基础设施及开发流程
加快软件开发和交付的速度
消除开发团队与运维团队之间的壁垒 



  
### SRE  
  
