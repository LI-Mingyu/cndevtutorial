- Role
- RoleBinding
- ServiceAccount

都是 namespace scope 的资源，如果创建时不注明 namespace ，均表示在执行命令的当前 namespace 创建。

如果在 RoleBinding 创建时 metadata 指定了 ns ，但是下面 subject 没有给 sa 指定 ns ，
那么即是对创建的 RoleBinding 所在的 ns 里的 sa 起作用，而不是对命令执行的当前 ns 里的 sa 起作用。

这里尤其具有迷惑性的是，default 这个 sa，哪里都有，但实际上不同 ns 的 default sa，是不同的实例。

从另一个角度看，Role 和 Rolebinding 也可以用来赋予 sa 跨 ns 的访问权限（而不用设置 ClusterRoleBinding ）。一个例子可以参考 RBACdemo2.yaml 。

---

- ClusterRole
- ClusterRoleBiniding 

都是 cluster scope 的资源，所以在 ClusterRoleBinding 创建时 subjets 里的 ServiceAccount 必须指明 namespace 。

---

关于 RBAC 的实验，可以使用 utils 里的 kubectl pod 配合测试。