import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../providers/providers.dart';
import '../utils/utils.dart';
import '../views/auth/login_screen.dart';
import '../views/auth/register_screen.dart';
import '../views/dashboard/dashboard_screen.dart';
import '../views/projects/project_list_screen.dart';
import '../views/projects/project_form_screen.dart';
import '../views/projects/project_detail_screen.dart';
import '../views/tasks/task_list_screen.dart';
import '../views/tasks/task_form_screen.dart';
import '../views/tasks/task_detail_screen.dart';

final _routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);

  return GoRouter(
    initialLocation: '/login',
    redirect: (context, state) {
      final isAuthenticated = authState.isAuthenticated;
      final isAuthRoute = state.matchedLocation == '/login' ||
          state.matchedLocation == '/register';

      if (!isAuthenticated && !isAuthRoute) {
        return '/login';
      }

      if (isAuthenticated && isAuthRoute) {
        return '/dashboard';
      }

      return null;
    },
    routes: [
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => const RegisterScreen(),
      ),
      GoRoute(
        path: '/dashboard',
        builder: (context, state) => const DashboardScreen(),
      ),
      GoRoute(
        path: '/projects',
        builder: (context, state) => const ProjectListScreen(),
        routes: [
          GoRoute(
            path: 'create',
            builder: (context, state) => const ProjectFormScreen(),
          ),
          GoRoute(
            path: ':id',
            builder: (context, state) {
              final id = int.parse(state.pathParameters['id']!);
              return ProjectDetailScreen(projectId: id);
            },
            routes: [
              GoRoute(
                path: 'edit',
                builder: (context, state) {
                  final id = int.parse(state.pathParameters['id']!);
                  return ProjectFormScreen(projectId: id);
                },
              ),
            ],
          ),
        ],
      ),
      GoRoute(
        path: '/tasks',
        builder: (context, state) => const TaskListScreen(),
        routes: [
          GoRoute(
            path: 'create',
            builder: (context, state) {
              final projectId = state.uri.queryParameters['project_id'];
              return TaskFormScreen(
                projectId: projectId != null ? int.tryParse(projectId) : null,
              );
            },
          ),
          GoRoute(
            path: ':id',
            builder: (context, state) {
              final id = int.parse(state.pathParameters['id']!);
              return TaskDetailScreen(taskId: id);
            },
            routes: [
              GoRoute(
                path: 'edit',
                builder: (context, state) {
                  final id = int.parse(state.pathParameters['id']!);
                  return TaskFormScreen(taskId: id);
                },
              ),
            ],
          ),
        ],
      ),
    ],
    errorBuilder: (context, state) => Scaffold(
      appBar: AppBar(title: const Text('错误')),
      body: Center(
        child: Text('页面不存在: ${state.matchedLocation}'),
      ),
    ),
  );
});

final sharedPreferencesFutureProvider = FutureProvider<SharedPreferences>((ref) {
  return SharedPreferences.getInstance();
});

class LifeNavigatorApp extends ConsumerWidget {
  const LifeNavigatorApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final sharedPreferencesAsync = ref.watch(sharedPreferencesFutureProvider);
    final router = ref.watch(_routerProvider);

    return sharedPreferencesAsync.when(
      loading: () => MaterialApp(
        home: Scaffold(
          body: Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const CircularProgressIndicator(),
                const SizedBox(height: 16),
                Text(
                  '咫尺生涯',
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
              ],
            ),
          ),
        ),
      ),
      error: (error, stack) => MaterialApp(
        home: Scaffold(
          body: Center(
            child: Text('初始化失败: $error'),
          ),
        ),
      ),
      data: (prefs) {
        return ProviderScope(
          overrides: [
            sharedPreferencesProvider.overrideWithValue(prefs),
          ],
          child: MaterialApp.router(
            title: '咫尺生涯',
            debugShowCheckedModeBanner: false,
            theme: AppTheme.lightTheme,
            darkTheme: AppTheme.darkTheme,
            themeMode: ThemeMode.system,
            routerConfig: router,
          ),
        );
      },
    );
  }
}
