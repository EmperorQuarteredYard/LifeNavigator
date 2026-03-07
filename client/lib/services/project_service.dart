import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/models.dart';
import 'api_client.dart';

class ProjectService {
  final ApiClient _client;

  ProjectService(this._client);

  Future<Project> createProject(CreateProjectRequest request) async {
    final response = await _client.dio.post(
      '/projects',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Project.fromJson(apiResponse.data!);
  }

  Future<Project> getProject(int id) async {
    final response = await _client.dio.get('/projects/$id');
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Project.fromJson(apiResponse.data!);
  }

  Future<ProjectListResponse> getProjects({
    int page = 1,
    int pageSize = 10,
  }) async {
    final response = await _client.dio.get(
      '/projects',
      queryParameters: {'page': page, 'page_size': pageSize},
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return ProjectListResponse.fromJson(apiResponse.data!);
  }

  Future<Project> updateProject(int id, UpdateProjectRequest request) async {
    final response = await _client.dio.put(
      '/projects/$id',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Project.fromJson(apiResponse.data!);
  }

  Future<void> deleteProject(int id) async {
    final response = await _client.dio.delete('/projects/$id');
    final apiResponse = ApiResponse<dynamic>.fromJson(
      response.data,
      (json) => json,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }

  Future<ProjectBudget> createBudget(
    int projectId,
    ProjectBudgetRequest request,
  ) async {
    final response = await _client.dio.post(
      '/projects/$projectId/budgets',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return ProjectBudget.fromJson(apiResponse.data!);
  }

  Future<ProjectBudget> updateBudget(
    int projectId,
    int budgetId,
    ProjectBudgetRequest request,
  ) async {
    final response = await _client.dio.put(
      '/projects/$projectId/budgets/$budgetId',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return ProjectBudget.fromJson(apiResponse.data!);
  }

  Future<void> deleteBudget(int projectId, int budgetId) async {
    final response = await _client.dio.delete(
      '/projects/$projectId/budgets/$budgetId',
    );
    final apiResponse = ApiResponse<dynamic>.fromJson(
      response.data,
      (json) => json,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }

  Future<ProjectBudgetSummary> getBudgetSummary(int projectId) async {
    final response = await _client.dio.get(
      '/projects/$projectId/budgets/summary',
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return ProjectBudgetSummary.fromJson(apiResponse.data!);
  }
}

final projectServiceProvider = Provider<ProjectService>((ref) {
  return ProjectService(ref.read(apiClientProvider));
});
