import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/models.dart';
import 'api_client.dart';

class TaskService {
  final ApiClient _client;

  TaskService(this._client);

  Future<Task> createTask(CreateTaskRequest request) async {
    final response = await _client.dio.post(
      '/tasks',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Task.fromJson(apiResponse.data!);
  }

  Future<Task> getTask(int id) async {
    final response = await _client.dio.get('/tasks/$id');
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Task.fromJson(apiResponse.data!);
  }

  Future<TaskListResponse> getTasks({
    int? projectId,
    int page = 1,
    int pageSize = 20,
    int? status,
    String? category,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (projectId != null) queryParams['project_id'] = projectId;
    if (status != null) queryParams['status'] = status;
    if (category != null) queryParams['category'] = category;

    final response = await _client.dio.get(
      '/tasks',
      queryParameters: queryParams,
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return TaskListResponse.fromJson(apiResponse.data!);
  }

  Future<Task> updateTask(int id, UpdateTaskRequest request) async {
    final response = await _client.dio.put(
      '/tasks/$id',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Task.fromJson(apiResponse.data!);
  }

  Future<void> deleteTask(int id) async {
    final response = await _client.dio.delete('/tasks/$id');
    final apiResponse = ApiResponse<dynamic>.fromJson(
      response.data,
      (json) => json,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }

  Future<Task> finishTask(int id, DateTime completedAt) async {
    final response = await _client.dio.post(
      '/tasks/$id/finish',
      data: FinishTaskRequest(time: completedAt.toIso8601String()).toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Task.fromJson(apiResponse.data!);
  }

  Future<TaskPayment> createPayment(
    int taskId,
    CreateTaskPaymentRequest request,
  ) async {
    final response = await _client.dio.post(
      '/tasks/$taskId/payments',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return TaskPayment.fromJson(apiResponse.data!);
  }

  Future<TaskPayment> updatePayment(
    int taskId,
    int paymentId,
    UpdateTaskPaymentRequest request,
  ) async {
    final response = await _client.dio.put(
      '/tasks/$taskId/payments/$paymentId',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return TaskPayment.fromJson(apiResponse.data!);
  }

  Future<void> deletePayment(int taskId, int paymentId) async {
    final response = await _client.dio.delete(
      '/tasks/$taskId/payments/$paymentId',
    );
    final apiResponse = ApiResponse<dynamic>.fromJson(
      response.data,
      (json) => json,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }

  Future<void> addPrerequisite(int taskId, int prerequisiteId) async {
    final response = await _client.dio.post(
      '/tasks/$taskId/prerequisites',
      data: PrerequisitesRequest(prerequisiteId: prerequisiteId).toJson(),
    );
    final apiResponse = ApiResponse<dynamic>.fromJson(
      response.data,
      (json) => json,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }

  Future<void> removePrerequisite(int taskId, int prerequisiteId) async {
    final response = await _client.dio.delete(
      '/tasks/$taskId/prerequisites/$prerequisiteId',
    );
    final apiResponse = ApiResponse<dynamic>.fromJson(
      response.data,
      (json) => json,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }

  Future<List<Dependency>> getPrerequisites(int taskId) async {
    final response = await _client.dio.get('/tasks/$taskId/prerequisites');
    final apiResponse = ApiResponse<List<dynamic>>.fromJson(
      response.data,
      (json) => json as List<dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!
        .map((e) => Dependency.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<List<Dependency>> getPostrequisites(int taskId) async {
    final response = await _client.dio.get('/tasks/$taskId/postrequisites');
    final apiResponse = ApiResponse<List<dynamic>>.fromJson(
      response.data,
      (json) => json as List<dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!
        .map((e) => Dependency.fromJson(e as Map<String, dynamic>))
        .toList();
  }
}

final taskServiceProvider = Provider<TaskService>((ref) {
  return TaskService(ref.read(apiClientProvider));
});
