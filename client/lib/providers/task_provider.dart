import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/models.dart';
import '../services/services.dart';

class TaskState {
  final List<Task> tasks;
  final Task? selectedTask;
  final bool isLoading;
  final String? error;
  final int currentPage;
  final int totalItems;
  final bool hasMore;
  final int? currentProjectId;

  const TaskState({
    this.tasks = const [],
    this.selectedTask,
    this.isLoading = false,
    this.error,
    this.currentPage = 1,
    this.totalItems = 0,
    this.hasMore = true,
    this.currentProjectId,
  });

  TaskState copyWith({
    List<Task>? tasks,
    Task? selectedTask,
    bool? isLoading,
    String? error,
    int? currentPage,
    int? totalItems,
    bool? hasMore,
    int? currentProjectId,
  }) {
    return TaskState(
      tasks: tasks ?? this.tasks,
      selectedTask: selectedTask ?? this.selectedTask,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      currentPage: currentPage ?? this.currentPage,
      totalItems: totalItems ?? this.totalItems,
      hasMore: hasMore ?? this.hasMore,
      currentProjectId: currentProjectId ?? this.currentProjectId,
    );
  }
}

class TaskNotifier extends StateNotifier<TaskState> {
  final TaskService _taskService;

  TaskNotifier(this._taskService) : super(const TaskState());

  Future<void> loadTasks({
    int? projectId,
    bool refresh = false,
    int? status,
    String? category,
  }) async {
    if (state.isLoading) return;

    if (refresh || projectId != state.currentProjectId) {
      state = TaskState(currentProjectId: projectId);
    }

    state = state.copyWith(isLoading: true, error: null);

    try {
      final page = refresh ? 1 : state.currentPage;
      final response = await _taskService.getTasks(
        projectId: projectId,
        page: page,
        status: status,
        category: category,
      );

      final newTasks = refresh
          ? response.list
          : [...state.tasks, ...response.list];

      state = state.copyWith(
        tasks: newTasks,
        isLoading: false,
        currentPage: page + 1,
        totalItems: response.total,
        hasMore: newTasks.length < response.total,
      );
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> selectTask(int id) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final task = await _taskService.getTask(id);
      state = state.copyWith(selectedTask: task, isLoading: false);
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<bool> createTask(CreateTaskRequest request) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final task = await _taskService.createTask(request);
      state = state.copyWith(
        tasks: [task, ...state.tasks],
        isLoading: false,
        totalItems: state.totalItems + 1,
      );
      return true;
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
      return false;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  Future<bool> updateTask(int id, UpdateTaskRequest request) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final updatedTask = await _taskService.updateTask(id, request);
      final newTasks = state.tasks.map((t) {
        return t.id == id ? updatedTask : t;
      }).toList();
      state = state.copyWith(
        tasks: newTasks,
        selectedTask: state.selectedTask?.id == id
            ? updatedTask
            : state.selectedTask,
        isLoading: false,
      );
      return true;
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
      return false;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  Future<bool> deleteTask(int id) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _taskService.deleteTask(id);
      final newTasks = state.tasks.where((t) => t.id != id).toList();
      state = state.copyWith(
        tasks: newTasks,
        selectedTask:
            state.selectedTask?.id == id ? null : state.selectedTask,
        isLoading: false,
        totalItems: state.totalItems - 1,
      );
      return true;
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
      return false;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  Future<bool> finishTask(int id, DateTime completedAt) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final updatedTask = await _taskService.finishTask(id, completedAt);
      final newTasks = state.tasks.map((t) {
        return t.id == id ? updatedTask : t;
      }).toList();
      state = state.copyWith(
        tasks: newTasks,
        selectedTask: state.selectedTask?.id == id
            ? updatedTask
            : state.selectedTask,
        isLoading: false,
      );
      return true;
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
      return false;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  Future<bool> createPayment(int taskId, CreateTaskPaymentRequest request) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final payment = await _taskService.createPayment(taskId, request);
      if (state.selectedTask?.id == taskId) {
        final updatedTask = state.selectedTask!.copyWith(
          payments: [...?state.selectedTask!.payments, payment],
        );
        state = state.copyWith(
          selectedTask: updatedTask,
          isLoading: false,
        );
      }
      return true;
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
      return false;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  void clearSelection() {
    state = state.copyWith(selectedTask: null);
  }
}

final taskStateProvider =
    StateNotifierProvider<TaskNotifier, TaskState>((ref) {
  return TaskNotifier(ref.read(taskServiceProvider));
});

final taskByIdProvider = Provider.family<Task?, int>((ref, id) {
  final state = ref.watch(taskStateProvider);
  return state.tasks.where((t) => t.id == id).firstOrNull;
});

final tasksByProjectProvider = Provider.family<List<Task>, int>((ref, projectId) {
  final state = ref.watch(taskStateProvider);
  return state.tasks.where((t) => t.projectId == projectId).toList();
});

final tasksByStatusProvider = Provider.family<List<Task>, TaskStatus>((ref, status) {
  final state = ref.watch(taskStateProvider);
  return state.tasks.where((t) => t.status == status.value).toList();
});
