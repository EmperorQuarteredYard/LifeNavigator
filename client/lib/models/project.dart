import 'package:freezed_annotation/freezed_annotation.dart';

part 'project.freezed.dart';
part 'project.g.dart';

@freezed
class ProjectBudget with _$ProjectBudget {
  const factory ProjectBudget({
    required int id,
    required int projectId,
    int? accountId,
    required double budget,
    required double used,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _ProjectBudget;

  factory ProjectBudget.fromJson(Map<String, dynamic> json) =>
      _$ProjectBudgetFromJson(json);
}

@freezed
class Project with _$Project {
  const factory Project({
    required int id,
    required int userId,
    required String name,
    String? description,
    required int refreshInterval,
    required DateTime lastRefresh,
    required int maxTaskId,
    required DateTime createdAt,
    required DateTime updatedAt,
    List<ProjectBudget>? budgets,
  }) = _Project;

  factory Project.fromJson(Map<String, dynamic> json) =>
      _$ProjectFromJson(json);
}

@freezed
class CreateProjectRequest with _$CreateProjectRequest {
  const factory CreateProjectRequest({
    required String name,
    String? description,
    int? refreshInterval,
  }) = _CreateProjectRequest;

  factory CreateProjectRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateProjectRequestFromJson(json);
}

@freezed
class UpdateProjectRequest with _$UpdateProjectRequest {
  const factory UpdateProjectRequest({
    required String name,
    String? description,
    int? refreshInterval,
  }) = _UpdateProjectRequest;

  factory UpdateProjectRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateProjectRequestFromJson(json);
}

@freezed
class ProjectBudgetRequest with _$ProjectBudgetRequest {
  const factory ProjectBudgetRequest({
    int? accountId,
    required double budget,
    double? used,
  }) = _ProjectBudgetRequest;

  factory ProjectBudgetRequest.fromJson(Map<String, dynamic> json) =>
      _$ProjectBudgetRequestFromJson(json);
}

@freezed
class ProjectBudgetSummary with _$ProjectBudgetSummary {
  const factory ProjectBudgetSummary({
    required List<ProjectBudget> budgets,
    required double totalBudget,
    required double totalUsed,
  }) = _ProjectBudgetSummary;

  factory ProjectBudgetSummary.fromJson(Map<String, dynamic> json) =>
      _$ProjectBudgetSummaryFromJson(json);
}

@freezed
class ProjectListResponse with _$ProjectListResponse {
  const factory ProjectListResponse({
    required int total,
    required List<Project> items,
  }) = _ProjectListResponse;

  factory ProjectListResponse.fromJson(Map<String, dynamic> json) =>
      _$ProjectListResponseFromJson(json);
}
