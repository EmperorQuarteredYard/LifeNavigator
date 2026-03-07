import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/models.dart';
import '../services/services.dart';

class ProjectState {
  final List<Project> projects;
  final Project? selectedProject;
  final bool isLoading;
  final String? error;
  final int currentPage;
  final int totalItems;
  final bool hasMore;

  const ProjectState({
    this.projects = const [],
    this.selectedProject,
    this.isLoading = false,
    this.error,
    this.currentPage = 1,
    this.totalItems = 0,
    this.hasMore = true,
  });

  ProjectState copyWith({
    List<Project>? projects,
    Project? selectedProject,
    bool? isLoading,
    String? error,
    int? currentPage,
    int? totalItems,
    bool? hasMore,
  }) {
    return ProjectState(
      projects: projects ?? this.projects,
      selectedProject: selectedProject ?? this.selectedProject,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      currentPage: currentPage ?? this.currentPage,
      totalItems: totalItems ?? this.totalItems,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

class ProjectNotifier extends StateNotifier<ProjectState> {
  final ProjectService _projectService;

  ProjectNotifier(this._projectService) : super(const ProjectState());

  Future<void> loadProjects({bool refresh = false}) async {
    if (state.isLoading) return;

    if (refresh) {
      state = const ProjectState();
    }

    state = state.copyWith(isLoading: true, error: null);

    try {
      final page = refresh ? 1 : state.currentPage;
      final response = await _projectService.getProjects(page: page);

      final newProjects = refresh
          ? response.items
          : [...state.projects, ...response.items];

      state = state.copyWith(
        projects: newProjects,
        isLoading: false,
        currentPage: page + 1,
        totalItems: response.total,
        hasMore: newProjects.length < response.total,
      );
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> selectProject(int id) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final project = await _projectService.getProject(id);
      state = state.copyWith(selectedProject: project, isLoading: false);
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<bool> createProject(CreateProjectRequest request) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final project = await _projectService.createProject(request);
      state = state.copyWith(
        projects: [project, ...state.projects],
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

  Future<bool> updateProject(int id, UpdateProjectRequest request) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final updatedProject = await _projectService.updateProject(id, request);
      final newProjects = state.projects.map((p) {
        return p.id == id ? updatedProject : p;
      }).toList();
      state = state.copyWith(
        projects: newProjects,
        selectedProject: state.selectedProject?.id == id
            ? updatedProject
            : state.selectedProject,
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

  Future<bool> deleteProject(int id) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _projectService.deleteProject(id);
      final newProjects =
          state.projects.where((p) => p.id != id).toList();
      state = state.copyWith(
        projects: newProjects,
        selectedProject:
            state.selectedProject?.id == id ? null : state.selectedProject,
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

  Future<bool> createBudget(int projectId, ProjectBudgetRequest request) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final budget = await _projectService.createBudget(projectId, request);
      if (state.selectedProject?.id == projectId) {
        final updatedProject = state.selectedProject!.copyWith(
          budgets: [...?state.selectedProject!.budgets, budget],
        );
        state = state.copyWith(
          selectedProject: updatedProject,
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
    state = state.copyWith(selectedProject: null);
  }
}

final projectStateProvider =
    StateNotifierProvider<ProjectNotifier, ProjectState>((ref) {
  return ProjectNotifier(ref.read(projectServiceProvider));
});

final projectByIdProvider = Provider.family<Project?, int>((ref, id) {
  final state = ref.watch(projectStateProvider);
  return state.projects.where((p) => p.id == id).firstOrNull;
});
