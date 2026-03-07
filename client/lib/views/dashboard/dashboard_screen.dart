import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/providers.dart';
import '../../widgets/widgets.dart';
import '../../utils/utils.dart';

class DashboardScreen extends ConsumerStatefulWidget {
  const DashboardScreen({super.key});

  @override
  ConsumerState<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends ConsumerState<DashboardScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(projectStateProvider.notifier).loadProjects(refresh: true);
      ref.read(taskStateProvider.notifier).loadTasks(refresh: true);
      ref.read(accountStateProvider.notifier).loadAccounts();
    });
  }

  @override
  Widget build(BuildContext context) {
    final projectState = ref.watch(projectStateProvider);
    final taskState = ref.watch(taskStateProvider);
    final accountState = ref.watch(accountStateProvider);
    final authState = ref.watch(authStateProvider);

    final pendingTasks = taskState.tasks
        .where((t) => t.status == TaskStatus.pending.value)
        .toList();
    final inProgressTasks = taskState.tasks
        .where((t) => t.status == TaskStatus.inProgress.value)
        .toList();
    final overdueTasks = taskState.tasks.where((t) {
      return t.status != TaskStatus.completed.value &&
          DateTimeUtils.isOverdue(t.deadline);
    }).toList();

    return Scaffold(
      appBar: AppBar(
        title: const Text('咫尺生涯'),
        actions: [
          IconButton(
            icon: const Icon(Icons.notifications_outlined),
            onPressed: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('通知功能开发中')),
              );
            },
          ),
          PopupMenuButton<String>(
            icon: CircleAvatar(
              radius: 16,
              child: Text(
                authState.user?.nickname ?? authState.user?.username ?? 'U',
                style: const TextStyle(fontSize: 12),
              ),
            ),
            onSelected: (value) {
              switch (value) {
                case 'profile':
                  context.push('/profile');
                  break;
                case 'settings':
                  context.push('/settings');
                  break;
                case 'logout':
                  ref.read(authStateProvider.notifier).logout();
                  context.go('/login');
                  break;
              }
            },
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'profile',
                child: ListTile(
                  leading: Icon(Icons.person_outline),
                  title: Text('个人资料'),
                ),
              ),
              const PopupMenuItem(
                value: 'settings',
                child: ListTile(
                  leading: Icon(Icons.settings_outlined),
                  title: Text('设置'),
                ),
              ),
              const PopupMenuItem(
                value: 'logout',
                child: ListTile(
                  leading: Icon(Icons.logout),
                  title: Text('退出登录'),
                ),
              ),
            ],
          ),
          const SizedBox(width: 8),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          await Future.wait([
            ref.read(projectStateProvider.notifier).loadProjects(refresh: true),
            ref.read(taskStateProvider.notifier).loadTasks(refresh: true),
            ref.read(accountStateProvider.notifier).loadAccounts(),
          ]);
        },
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildWelcomeCard(authState),
              const SizedBox(height: 16),
              _buildQuickStats(
                context,
                pendingTasks.length,
                inProgressTasks.length,
                overdueTasks.length,
              ),
              const SizedBox(height: 24),
              _buildSectionTitle(context, '账户概览', Icons.account_balance_wallet),
              const SizedBox(height: 12),
              _buildAccountsSummary(accountState),
              const SizedBox(height: 24),
              _buildSectionTitle(context, '今日待办', Icons.today),
              const SizedBox(height: 12),
              if (taskState.isLoading)
                const Center(child: CircularProgressIndicator())
              else if (pendingTasks.isEmpty && inProgressTasks.isEmpty)
                const EmptyState(
                  title: '暂无待办任务',
                  subtitle: '点击右下角按钮创建新任务',
                  icon: Icons.check_circle_outline,
                )
              else
                ...inProgressTasks.take(3).map((task) => TaskCard(
                      task: task,
                      onTap: () => context.push('/tasks/${task.id}'),
                      onComplete: () async {
                        await ref
                            .read(taskStateProvider.notifier)
                            .finishTask(task.id, DateTime.now());
                      },
                    )),
              ...pendingTasks.take(5).map((task) => TaskCard(
                    task: task,
                    onTap: () => context.push('/tasks/${task.id}'),
                    onComplete: () async {
                      await ref
                          .read(taskStateProvider.notifier)
                          .finishTask(task.id, DateTime.now());
                    },
                  )),
              const SizedBox(height: 24),
              _buildSectionTitle(context, '我的项目', Icons.folder_outlined),
              const SizedBox(height: 12),
              if (projectState.isLoading)
                const Center(child: CircularProgressIndicator())
              else if (projectState.projects.isEmpty)
                const EmptyState(
                  title: '暂无项目',
                  subtitle: '点击右下角按钮创建新项目',
                  icon: Icons.folder_outlined,
                )
              else
                ...projectState.projects.take(3).map((project) => ProjectCard(
                      project: project,
                      onTap: () => context.push('/projects/${project.id}'),
                    )),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          showModalBottomSheet(
            context: context,
            builder: (context) => _buildQuickActions(context),
          );
        },
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildWelcomeCard(AuthState authState) {
    final hour = DateTime.now().hour;
    String greeting;
    if (hour < 6) {
      greeting = '夜深了';
    } else if (hour < 12) {
      greeting = '早上好';
    } else if (hour < 14) {
      greeting = '中午好';
    } else if (hour < 18) {
      greeting = '下午好';
    } else {
      greeting = '晚上好';
    }

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '$greeting，${authState.user?.nickname ?? authState.user?.username ?? '用户'}',
                    style: Theme.of(context).textTheme.titleLarge?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    DateTimeUtils.formatDate(DateTime.now()),
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: Colors.grey.shade600,
                        ),
                  ),
                ],
              ),
            ),
            const Icon(
              Icons.wb_sunny_outlined,
              size: 48,
              color: Colors.orange,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildQuickStats(
    BuildContext context,
    int pendingCount,
    int inProgressCount,
    int overdueCount,
  ) {
    return Row(
      children: [
        Expanded(
          child: _StatCard(
            title: '待处理',
            value: pendingCount.toString(),
            icon: Icons.pending_actions_outlined,
            color: AppTheme.pendingColor,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _StatCard(
            title: '进行中',
            value: inProgressCount.toString(),
            icon: Icons.autorenew,
            color: AppTheme.inProgressColor,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _StatCard(
            title: '已过期',
            value: overdueCount.toString(),
            icon: Icons.warning_amber_outlined,
            color: AppTheme.error,
          ),
        ),
      ],
    );
  }

  Widget _buildAccountsSummary(AccountState accountState) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '总资产',
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: Colors.grey.shade600,
                      ),
                ),
                Text(
                  CurrencyUtils.format(accountState.totalBalance),
                  style: Theme.of(context).textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: AppTheme.success,
                      ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            const Divider(),
            const SizedBox(height: 8),
            if (accountState.accounts.isEmpty)
              Text(
                '暂无账户',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Colors.grey.shade500,
                    ),
              )
            else
              ...accountState.accounts.take(3).map((account) => Padding(
                    padding: const EdgeInsets.symmetric(vertical: 4),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          children: [
                            Icon(
                              _getAccountIcon(account.type),
                              size: 20,
                              color: Colors.grey.shade600,
                            ),
                            const SizedBox(width: 8),
                            Text(account.type),
                          ],
                        ),
                        Text(
                          CurrencyUtils.format(account.balance),
                          style: Theme.of(context).textTheme.bodyMedium,
                        ),
                      ],
                    ),
                  )),
          ],
        ),
      ),
    );
  }

  IconData _getAccountIcon(String type) {
    switch (type.toLowerCase()) {
      case 'cash':
      case '现金':
        return Icons.money_outlined;
      case 'bank':
      case '银行卡':
        return Icons.account_balance_outlined;
      case 'credit':
      case '信用卡':
        return Icons.credit_card_outlined;
      case 'alipay':
      case '支付宝':
        return Icons.payment_outlined;
      case 'wechat':
      case '微信':
        return Icons.chat_outlined;
      default:
        return Icons.account_balance_wallet_outlined;
    }
  }

  Widget _buildSectionTitle(BuildContext context, String title, IconData icon) {
    return Row(
      children: [
        Icon(icon, size: 20, color: Theme.of(context).colorScheme.primary),
        const SizedBox(width: 8),
        Text(
          title,
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
        ),
      ],
    );
  }

  Widget _buildQuickActions(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ListTile(
            leading: const Icon(Icons.folder_outlined),
            title: const Text('新建项目'),
            onTap: () {
              Navigator.pop(context);
              context.push('/projects/create');
            },
          ),
          ListTile(
            leading: const Icon(Icons.task_outlined),
            title: const Text('新建任务'),
            onTap: () {
              Navigator.pop(context);
              context.push('/tasks/create');
            },
          ),
          ListTile(
            leading: const Icon(Icons.account_balance_wallet_outlined),
            title: const Text('新建账户'),
            onTap: () {
              Navigator.pop(context);
              context.push('/accounts/create');
            },
          ),
        ],
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;
  final Color color;

  const _StatCard({
    required this.title,
    required this.value,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Icon(icon, color: color, size: 28),
            const SizedBox(height: 8),
            Text(
              value,
              style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: color,
                  ),
            ),
            const SizedBox(height: 4),
            Text(
              title,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: Colors.grey.shade600,
                  ),
            ),
          ],
        ),
      ),
    );
  }
}
