import 'package:freezed_annotation/freezed_annotation.dart';

part 'task.freezed.dart';
part 'task.g.dart';

@freezed
class TaskPayment with _$TaskPayment {
  const factory TaskPayment({
    required int id,
    required int taskId,
    required int budgetId,
    required double amount,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _TaskPayment;

  factory TaskPayment.fromJson(Map<String, dynamic> json) =>
      _$TaskPaymentFromJson(json);
}

@freezed
class Task with _$Task {
  const factory Task({
    required int id,
    required int userId,
    required int projectId,
    required String name,
    String? description,
    required bool autoCalculated,
    required int type,
    required int status,
    String? category,
    DateTime? deadline,
    DateTime? completedAt,
    required DateTime createdAt,
    required DateTime updatedAt,
    List<TaskPayment>? payments,
  }) = _Task;

  factory Task.fromJson(Map<String, dynamic> json) => _$TaskFromJson(json);
}

@freezed
class CreateTaskRequest with _$CreateTaskRequest {
  const factory CreateTaskRequest({
    required int projectId,
    required String name,
    String? description,
    bool? autoCalculated,
    int? type,
    int? status,
    String? category,
    String? deadline,
  }) = _CreateTaskRequest;

  factory CreateTaskRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateTaskRequestFromJson(json);
}

@freezed
class UpdateTaskRequest with _$UpdateTaskRequest {
  const factory UpdateTaskRequest({
    String? name,
    int? projectId,
    String? description,
    bool? autoCalculated,
    int? type,
    int? status,
    String? category,
    String? deadline,
  }) = _UpdateTaskRequest;

  factory UpdateTaskRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateTaskRequestFromJson(json);
}

@freezed
class FinishTaskRequest with _$FinishTaskRequest {
  const factory FinishTaskRequest({
    required String time,
  }) = _FinishTaskRequest;

  factory FinishTaskRequest.fromJson(Map<String, dynamic> json) =>
      _$FinishTaskRequestFromJson(json);
}

@freezed
class CreateTaskPaymentRequest with _$CreateTaskPaymentRequest {
  const factory CreateTaskPaymentRequest({
    required int budgetId,
    required double amount,
  }) = _CreateTaskPaymentRequest;

  factory CreateTaskPaymentRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateTaskPaymentRequestFromJson(json);
}

@freezed
class UpdateTaskPaymentRequest with _$UpdateTaskPaymentRequest {
  const factory UpdateTaskPaymentRequest({
    required double amount,
  }) = _UpdateTaskPaymentRequest;

  factory UpdateTaskPaymentRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateTaskPaymentRequestFromJson(json);
}

@freezed
class PrerequisitesRequest with _$PrerequisitesRequest {
  const factory PrerequisitesRequest({
    required int prerequisiteId,
  }) = _PrerequisitesRequest;

  factory PrerequisitesRequest.fromJson(Map<String, dynamic> json) =>
      _$PrerequisitesRequestFromJson(json);
}

@freezed
class Dependency with _$Dependency {
  const factory Dependency({
    required int prerequisiteId,
    required int taskId,
  }) = _Dependency;

  factory Dependency.fromJson(Map<String, dynamic> json) =>
      _$DependencyFromJson(json);
}

@freezed
class TaskListResponse with _$TaskListResponse {
  const factory TaskListResponse({
    required int page,
    required int total,
    required int size,
    required List<Task> list,
  }) = _TaskListResponse;

  factory TaskListResponse.fromJson(Map<String, dynamic> json) =>
      _$TaskListResponseFromJson(json);
}

enum TaskStatus {
  pending(0),
  inProgress(1),
  completed(2),
  cancelled(3);

  const TaskStatus(this.value);
  final int value;

  static TaskStatus fromValue(int value) {
    return TaskStatus.values.firstWhere(
      (e) => e.value == value,
      orElse: () => TaskStatus.pending,
    );
  }
}

enum TaskType {
  normal(0),
  milestone(1),
  routine(2);

  const TaskType(this.value);
  final int value;

  static TaskType fromValue(int value) {
    return TaskType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => TaskType.normal,
    );
  }
}
