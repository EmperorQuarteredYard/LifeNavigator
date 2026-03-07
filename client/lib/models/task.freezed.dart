// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'task.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

TaskPayment _$TaskPaymentFromJson(Map<String, dynamic> json) {
  return _TaskPayment.fromJson(json);
}

/// @nodoc
mixin _$TaskPayment {
  int get id => throw _privateConstructorUsedError;
  int get taskId => throw _privateConstructorUsedError;
  int get budgetId => throw _privateConstructorUsedError;
  double get amount => throw _privateConstructorUsedError;
  DateTime get createdAt => throw _privateConstructorUsedError;
  DateTime get updatedAt => throw _privateConstructorUsedError;

  /// Serializes this TaskPayment to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of TaskPayment
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $TaskPaymentCopyWith<TaskPayment> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $TaskPaymentCopyWith<$Res> {
  factory $TaskPaymentCopyWith(
          TaskPayment value, $Res Function(TaskPayment) then) =
      _$TaskPaymentCopyWithImpl<$Res, TaskPayment>;
  @useResult
  $Res call(
      {int id,
      int taskId,
      int budgetId,
      double amount,
      DateTime createdAt,
      DateTime updatedAt});
}

/// @nodoc
class _$TaskPaymentCopyWithImpl<$Res, $Val extends TaskPayment>
    implements $TaskPaymentCopyWith<$Res> {
  _$TaskPaymentCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of TaskPayment
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? taskId = null,
    Object? budgetId = null,
    Object? amount = null,
    Object? createdAt = null,
    Object? updatedAt = null,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      taskId: null == taskId
          ? _value.taskId
          : taskId // ignore: cast_nullable_to_non_nullable
              as int,
      budgetId: null == budgetId
          ? _value.budgetId
          : budgetId // ignore: cast_nullable_to_non_nullable
              as int,
      amount: null == amount
          ? _value.amount
          : amount // ignore: cast_nullable_to_non_nullable
              as double,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$TaskPaymentImplCopyWith<$Res>
    implements $TaskPaymentCopyWith<$Res> {
  factory _$$TaskPaymentImplCopyWith(
          _$TaskPaymentImpl value, $Res Function(_$TaskPaymentImpl) then) =
      __$$TaskPaymentImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      int taskId,
      int budgetId,
      double amount,
      DateTime createdAt,
      DateTime updatedAt});
}

/// @nodoc
class __$$TaskPaymentImplCopyWithImpl<$Res>
    extends _$TaskPaymentCopyWithImpl<$Res, _$TaskPaymentImpl>
    implements _$$TaskPaymentImplCopyWith<$Res> {
  __$$TaskPaymentImplCopyWithImpl(
      _$TaskPaymentImpl _value, $Res Function(_$TaskPaymentImpl) _then)
      : super(_value, _then);

  /// Create a copy of TaskPayment
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? taskId = null,
    Object? budgetId = null,
    Object? amount = null,
    Object? createdAt = null,
    Object? updatedAt = null,
  }) {
    return _then(_$TaskPaymentImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      taskId: null == taskId
          ? _value.taskId
          : taskId // ignore: cast_nullable_to_non_nullable
              as int,
      budgetId: null == budgetId
          ? _value.budgetId
          : budgetId // ignore: cast_nullable_to_non_nullable
              as int,
      amount: null == amount
          ? _value.amount
          : amount // ignore: cast_nullable_to_non_nullable
              as double,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$TaskPaymentImpl implements _TaskPayment {
  const _$TaskPaymentImpl(
      {required this.id,
      required this.taskId,
      required this.budgetId,
      required this.amount,
      required this.createdAt,
      required this.updatedAt});

  factory _$TaskPaymentImpl.fromJson(Map<String, dynamic> json) =>
      _$$TaskPaymentImplFromJson(json);

  @override
  final int id;
  @override
  final int taskId;
  @override
  final int budgetId;
  @override
  final double amount;
  @override
  final DateTime createdAt;
  @override
  final DateTime updatedAt;

  @override
  String toString() {
    return 'TaskPayment(id: $id, taskId: $taskId, budgetId: $budgetId, amount: $amount, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$TaskPaymentImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.taskId, taskId) || other.taskId == taskId) &&
            (identical(other.budgetId, budgetId) ||
                other.budgetId == budgetId) &&
            (identical(other.amount, amount) || other.amount == amount) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType, id, taskId, budgetId, amount, createdAt, updatedAt);

  /// Create a copy of TaskPayment
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$TaskPaymentImplCopyWith<_$TaskPaymentImpl> get copyWith =>
      __$$TaskPaymentImplCopyWithImpl<_$TaskPaymentImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$TaskPaymentImplToJson(
      this,
    );
  }
}

abstract class _TaskPayment implements TaskPayment {
  const factory _TaskPayment(
      {required final int id,
      required final int taskId,
      required final int budgetId,
      required final double amount,
      required final DateTime createdAt,
      required final DateTime updatedAt}) = _$TaskPaymentImpl;

  factory _TaskPayment.fromJson(Map<String, dynamic> json) =
      _$TaskPaymentImpl.fromJson;

  @override
  int get id;
  @override
  int get taskId;
  @override
  int get budgetId;
  @override
  double get amount;
  @override
  DateTime get createdAt;
  @override
  DateTime get updatedAt;

  /// Create a copy of TaskPayment
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$TaskPaymentImplCopyWith<_$TaskPaymentImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

Task _$TaskFromJson(Map<String, dynamic> json) {
  return _Task.fromJson(json);
}

/// @nodoc
mixin _$Task {
  int get id => throw _privateConstructorUsedError;
  int get userId => throw _privateConstructorUsedError;
  int get projectId => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  String? get description => throw _privateConstructorUsedError;
  bool get autoCalculated => throw _privateConstructorUsedError;
  int get type => throw _privateConstructorUsedError;
  int get status => throw _privateConstructorUsedError;
  String? get category => throw _privateConstructorUsedError;
  DateTime? get deadline => throw _privateConstructorUsedError;
  DateTime? get completedAt => throw _privateConstructorUsedError;
  DateTime get createdAt => throw _privateConstructorUsedError;
  DateTime get updatedAt => throw _privateConstructorUsedError;
  List<TaskPayment>? get payments => throw _privateConstructorUsedError;

  /// Serializes this Task to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Task
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $TaskCopyWith<Task> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $TaskCopyWith<$Res> {
  factory $TaskCopyWith(Task value, $Res Function(Task) then) =
      _$TaskCopyWithImpl<$Res, Task>;
  @useResult
  $Res call(
      {int id,
      int userId,
      int projectId,
      String name,
      String? description,
      bool autoCalculated,
      int type,
      int status,
      String? category,
      DateTime? deadline,
      DateTime? completedAt,
      DateTime createdAt,
      DateTime updatedAt,
      List<TaskPayment>? payments});
}

/// @nodoc
class _$TaskCopyWithImpl<$Res, $Val extends Task>
    implements $TaskCopyWith<$Res> {
  _$TaskCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Task
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? projectId = null,
    Object? name = null,
    Object? description = freezed,
    Object? autoCalculated = null,
    Object? type = null,
    Object? status = null,
    Object? category = freezed,
    Object? deadline = freezed,
    Object? completedAt = freezed,
    Object? createdAt = null,
    Object? updatedAt = null,
    Object? payments = freezed,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      userId: null == userId
          ? _value.userId
          : userId // ignore: cast_nullable_to_non_nullable
              as int,
      projectId: null == projectId
          ? _value.projectId
          : projectId // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      autoCalculated: null == autoCalculated
          ? _value.autoCalculated
          : autoCalculated // ignore: cast_nullable_to_non_nullable
              as bool,
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as int,
      status: null == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int,
      category: freezed == category
          ? _value.category
          : category // ignore: cast_nullable_to_non_nullable
              as String?,
      deadline: freezed == deadline
          ? _value.deadline
          : deadline // ignore: cast_nullable_to_non_nullable
              as DateTime?,
      completedAt: freezed == completedAt
          ? _value.completedAt
          : completedAt // ignore: cast_nullable_to_non_nullable
              as DateTime?,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      payments: freezed == payments
          ? _value.payments
          : payments // ignore: cast_nullable_to_non_nullable
              as List<TaskPayment>?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$TaskImplCopyWith<$Res> implements $TaskCopyWith<$Res> {
  factory _$$TaskImplCopyWith(
          _$TaskImpl value, $Res Function(_$TaskImpl) then) =
      __$$TaskImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      int userId,
      int projectId,
      String name,
      String? description,
      bool autoCalculated,
      int type,
      int status,
      String? category,
      DateTime? deadline,
      DateTime? completedAt,
      DateTime createdAt,
      DateTime updatedAt,
      List<TaskPayment>? payments});
}

/// @nodoc
class __$$TaskImplCopyWithImpl<$Res>
    extends _$TaskCopyWithImpl<$Res, _$TaskImpl>
    implements _$$TaskImplCopyWith<$Res> {
  __$$TaskImplCopyWithImpl(_$TaskImpl _value, $Res Function(_$TaskImpl) _then)
      : super(_value, _then);

  /// Create a copy of Task
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? projectId = null,
    Object? name = null,
    Object? description = freezed,
    Object? autoCalculated = null,
    Object? type = null,
    Object? status = null,
    Object? category = freezed,
    Object? deadline = freezed,
    Object? completedAt = freezed,
    Object? createdAt = null,
    Object? updatedAt = null,
    Object? payments = freezed,
  }) {
    return _then(_$TaskImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      userId: null == userId
          ? _value.userId
          : userId // ignore: cast_nullable_to_non_nullable
              as int,
      projectId: null == projectId
          ? _value.projectId
          : projectId // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      autoCalculated: null == autoCalculated
          ? _value.autoCalculated
          : autoCalculated // ignore: cast_nullable_to_non_nullable
              as bool,
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as int,
      status: null == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int,
      category: freezed == category
          ? _value.category
          : category // ignore: cast_nullable_to_non_nullable
              as String?,
      deadline: freezed == deadline
          ? _value.deadline
          : deadline // ignore: cast_nullable_to_non_nullable
              as DateTime?,
      completedAt: freezed == completedAt
          ? _value.completedAt
          : completedAt // ignore: cast_nullable_to_non_nullable
              as DateTime?,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      payments: freezed == payments
          ? _value._payments
          : payments // ignore: cast_nullable_to_non_nullable
              as List<TaskPayment>?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$TaskImpl implements _Task {
  const _$TaskImpl(
      {required this.id,
      required this.userId,
      required this.projectId,
      required this.name,
      this.description,
      required this.autoCalculated,
      required this.type,
      required this.status,
      this.category,
      this.deadline,
      this.completedAt,
      required this.createdAt,
      required this.updatedAt,
      final List<TaskPayment>? payments})
      : _payments = payments;

  factory _$TaskImpl.fromJson(Map<String, dynamic> json) =>
      _$$TaskImplFromJson(json);

  @override
  final int id;
  @override
  final int userId;
  @override
  final int projectId;
  @override
  final String name;
  @override
  final String? description;
  @override
  final bool autoCalculated;
  @override
  final int type;
  @override
  final int status;
  @override
  final String? category;
  @override
  final DateTime? deadline;
  @override
  final DateTime? completedAt;
  @override
  final DateTime createdAt;
  @override
  final DateTime updatedAt;
  final List<TaskPayment>? _payments;
  @override
  List<TaskPayment>? get payments {
    final value = _payments;
    if (value == null) return null;
    if (_payments is EqualUnmodifiableListView) return _payments;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(value);
  }

  @override
  String toString() {
    return 'Task(id: $id, userId: $userId, projectId: $projectId, name: $name, description: $description, autoCalculated: $autoCalculated, type: $type, status: $status, category: $category, deadline: $deadline, completedAt: $completedAt, createdAt: $createdAt, updatedAt: $updatedAt, payments: $payments)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$TaskImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.userId, userId) || other.userId == userId) &&
            (identical(other.projectId, projectId) ||
                other.projectId == projectId) &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.autoCalculated, autoCalculated) ||
                other.autoCalculated == autoCalculated) &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.category, category) ||
                other.category == category) &&
            (identical(other.deadline, deadline) ||
                other.deadline == deadline) &&
            (identical(other.completedAt, completedAt) ||
                other.completedAt == completedAt) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt) &&
            const DeepCollectionEquality().equals(other._payments, _payments));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType,
      id,
      userId,
      projectId,
      name,
      description,
      autoCalculated,
      type,
      status,
      category,
      deadline,
      completedAt,
      createdAt,
      updatedAt,
      const DeepCollectionEquality().hash(_payments));

  /// Create a copy of Task
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$TaskImplCopyWith<_$TaskImpl> get copyWith =>
      __$$TaskImplCopyWithImpl<_$TaskImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$TaskImplToJson(
      this,
    );
  }
}

abstract class _Task implements Task {
  const factory _Task(
      {required final int id,
      required final int userId,
      required final int projectId,
      required final String name,
      final String? description,
      required final bool autoCalculated,
      required final int type,
      required final int status,
      final String? category,
      final DateTime? deadline,
      final DateTime? completedAt,
      required final DateTime createdAt,
      required final DateTime updatedAt,
      final List<TaskPayment>? payments}) = _$TaskImpl;

  factory _Task.fromJson(Map<String, dynamic> json) = _$TaskImpl.fromJson;

  @override
  int get id;
  @override
  int get userId;
  @override
  int get projectId;
  @override
  String get name;
  @override
  String? get description;
  @override
  bool get autoCalculated;
  @override
  int get type;
  @override
  int get status;
  @override
  String? get category;
  @override
  DateTime? get deadline;
  @override
  DateTime? get completedAt;
  @override
  DateTime get createdAt;
  @override
  DateTime get updatedAt;
  @override
  List<TaskPayment>? get payments;

  /// Create a copy of Task
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$TaskImplCopyWith<_$TaskImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateTaskRequest _$CreateTaskRequestFromJson(Map<String, dynamic> json) {
  return _CreateTaskRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateTaskRequest {
  int get projectId => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  String? get description => throw _privateConstructorUsedError;
  bool? get autoCalculated => throw _privateConstructorUsedError;
  int? get type => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError;
  String? get category => throw _privateConstructorUsedError;
  String? get deadline => throw _privateConstructorUsedError;

  /// Serializes this CreateTaskRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of CreateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CreateTaskRequestCopyWith<CreateTaskRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateTaskRequestCopyWith<$Res> {
  factory $CreateTaskRequestCopyWith(
          CreateTaskRequest value, $Res Function(CreateTaskRequest) then) =
      _$CreateTaskRequestCopyWithImpl<$Res, CreateTaskRequest>;
  @useResult
  $Res call(
      {int projectId,
      String name,
      String? description,
      bool? autoCalculated,
      int? type,
      int? status,
      String? category,
      String? deadline});
}

/// @nodoc
class _$CreateTaskRequestCopyWithImpl<$Res, $Val extends CreateTaskRequest>
    implements $CreateTaskRequestCopyWith<$Res> {
  _$CreateTaskRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CreateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? projectId = null,
    Object? name = null,
    Object? description = freezed,
    Object? autoCalculated = freezed,
    Object? type = freezed,
    Object? status = freezed,
    Object? category = freezed,
    Object? deadline = freezed,
  }) {
    return _then(_value.copyWith(
      projectId: null == projectId
          ? _value.projectId
          : projectId // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      autoCalculated: freezed == autoCalculated
          ? _value.autoCalculated
          : autoCalculated // ignore: cast_nullable_to_non_nullable
              as bool?,
      type: freezed == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      category: freezed == category
          ? _value.category
          : category // ignore: cast_nullable_to_non_nullable
              as String?,
      deadline: freezed == deadline
          ? _value.deadline
          : deadline // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CreateTaskRequestImplCopyWith<$Res>
    implements $CreateTaskRequestCopyWith<$Res> {
  factory _$$CreateTaskRequestImplCopyWith(_$CreateTaskRequestImpl value,
          $Res Function(_$CreateTaskRequestImpl) then) =
      __$$CreateTaskRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int projectId,
      String name,
      String? description,
      bool? autoCalculated,
      int? type,
      int? status,
      String? category,
      String? deadline});
}

/// @nodoc
class __$$CreateTaskRequestImplCopyWithImpl<$Res>
    extends _$CreateTaskRequestCopyWithImpl<$Res, _$CreateTaskRequestImpl>
    implements _$$CreateTaskRequestImplCopyWith<$Res> {
  __$$CreateTaskRequestImplCopyWithImpl(_$CreateTaskRequestImpl _value,
      $Res Function(_$CreateTaskRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of CreateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? projectId = null,
    Object? name = null,
    Object? description = freezed,
    Object? autoCalculated = freezed,
    Object? type = freezed,
    Object? status = freezed,
    Object? category = freezed,
    Object? deadline = freezed,
  }) {
    return _then(_$CreateTaskRequestImpl(
      projectId: null == projectId
          ? _value.projectId
          : projectId // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      autoCalculated: freezed == autoCalculated
          ? _value.autoCalculated
          : autoCalculated // ignore: cast_nullable_to_non_nullable
              as bool?,
      type: freezed == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      category: freezed == category
          ? _value.category
          : category // ignore: cast_nullable_to_non_nullable
              as String?,
      deadline: freezed == deadline
          ? _value.deadline
          : deadline // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CreateTaskRequestImpl implements _CreateTaskRequest {
  const _$CreateTaskRequestImpl(
      {required this.projectId,
      required this.name,
      this.description,
      this.autoCalculated,
      this.type,
      this.status,
      this.category,
      this.deadline});

  factory _$CreateTaskRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$CreateTaskRequestImplFromJson(json);

  @override
  final int projectId;
  @override
  final String name;
  @override
  final String? description;
  @override
  final bool? autoCalculated;
  @override
  final int? type;
  @override
  final int? status;
  @override
  final String? category;
  @override
  final String? deadline;

  @override
  String toString() {
    return 'CreateTaskRequest(projectId: $projectId, name: $name, description: $description, autoCalculated: $autoCalculated, type: $type, status: $status, category: $category, deadline: $deadline)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateTaskRequestImpl &&
            (identical(other.projectId, projectId) ||
                other.projectId == projectId) &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.autoCalculated, autoCalculated) ||
                other.autoCalculated == autoCalculated) &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.category, category) ||
                other.category == category) &&
            (identical(other.deadline, deadline) ||
                other.deadline == deadline));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, projectId, name, description,
      autoCalculated, type, status, category, deadline);

  /// Create a copy of CreateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateTaskRequestImplCopyWith<_$CreateTaskRequestImpl> get copyWith =>
      __$$CreateTaskRequestImplCopyWithImpl<_$CreateTaskRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateTaskRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateTaskRequest implements CreateTaskRequest {
  const factory _CreateTaskRequest(
      {required final int projectId,
      required final String name,
      final String? description,
      final bool? autoCalculated,
      final int? type,
      final int? status,
      final String? category,
      final String? deadline}) = _$CreateTaskRequestImpl;

  factory _CreateTaskRequest.fromJson(Map<String, dynamic> json) =
      _$CreateTaskRequestImpl.fromJson;

  @override
  int get projectId;
  @override
  String get name;
  @override
  String? get description;
  @override
  bool? get autoCalculated;
  @override
  int? get type;
  @override
  int? get status;
  @override
  String? get category;
  @override
  String? get deadline;

  /// Create a copy of CreateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CreateTaskRequestImplCopyWith<_$CreateTaskRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

UpdateTaskRequest _$UpdateTaskRequestFromJson(Map<String, dynamic> json) {
  return _UpdateTaskRequest.fromJson(json);
}

/// @nodoc
mixin _$UpdateTaskRequest {
  String? get name => throw _privateConstructorUsedError;
  int? get projectId => throw _privateConstructorUsedError;
  String? get description => throw _privateConstructorUsedError;
  bool? get autoCalculated => throw _privateConstructorUsedError;
  int? get type => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError;
  String? get category => throw _privateConstructorUsedError;
  String? get deadline => throw _privateConstructorUsedError;

  /// Serializes this UpdateTaskRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of UpdateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UpdateTaskRequestCopyWith<UpdateTaskRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UpdateTaskRequestCopyWith<$Res> {
  factory $UpdateTaskRequestCopyWith(
          UpdateTaskRequest value, $Res Function(UpdateTaskRequest) then) =
      _$UpdateTaskRequestCopyWithImpl<$Res, UpdateTaskRequest>;
  @useResult
  $Res call(
      {String? name,
      int? projectId,
      String? description,
      bool? autoCalculated,
      int? type,
      int? status,
      String? category,
      String? deadline});
}

/// @nodoc
class _$UpdateTaskRequestCopyWithImpl<$Res, $Val extends UpdateTaskRequest>
    implements $UpdateTaskRequestCopyWith<$Res> {
  _$UpdateTaskRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of UpdateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? projectId = freezed,
    Object? description = freezed,
    Object? autoCalculated = freezed,
    Object? type = freezed,
    Object? status = freezed,
    Object? category = freezed,
    Object? deadline = freezed,
  }) {
    return _then(_value.copyWith(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      projectId: freezed == projectId
          ? _value.projectId
          : projectId // ignore: cast_nullable_to_non_nullable
              as int?,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      autoCalculated: freezed == autoCalculated
          ? _value.autoCalculated
          : autoCalculated // ignore: cast_nullable_to_non_nullable
              as bool?,
      type: freezed == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      category: freezed == category
          ? _value.category
          : category // ignore: cast_nullable_to_non_nullable
              as String?,
      deadline: freezed == deadline
          ? _value.deadline
          : deadline // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$UpdateTaskRequestImplCopyWith<$Res>
    implements $UpdateTaskRequestCopyWith<$Res> {
  factory _$$UpdateTaskRequestImplCopyWith(_$UpdateTaskRequestImpl value,
          $Res Function(_$UpdateTaskRequestImpl) then) =
      __$$UpdateTaskRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String? name,
      int? projectId,
      String? description,
      bool? autoCalculated,
      int? type,
      int? status,
      String? category,
      String? deadline});
}

/// @nodoc
class __$$UpdateTaskRequestImplCopyWithImpl<$Res>
    extends _$UpdateTaskRequestCopyWithImpl<$Res, _$UpdateTaskRequestImpl>
    implements _$$UpdateTaskRequestImplCopyWith<$Res> {
  __$$UpdateTaskRequestImplCopyWithImpl(_$UpdateTaskRequestImpl _value,
      $Res Function(_$UpdateTaskRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of UpdateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? projectId = freezed,
    Object? description = freezed,
    Object? autoCalculated = freezed,
    Object? type = freezed,
    Object? status = freezed,
    Object? category = freezed,
    Object? deadline = freezed,
  }) {
    return _then(_$UpdateTaskRequestImpl(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      projectId: freezed == projectId
          ? _value.projectId
          : projectId // ignore: cast_nullable_to_non_nullable
              as int?,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      autoCalculated: freezed == autoCalculated
          ? _value.autoCalculated
          : autoCalculated // ignore: cast_nullable_to_non_nullable
              as bool?,
      type: freezed == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      category: freezed == category
          ? _value.category
          : category // ignore: cast_nullable_to_non_nullable
              as String?,
      deadline: freezed == deadline
          ? _value.deadline
          : deadline // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UpdateTaskRequestImpl implements _UpdateTaskRequest {
  const _$UpdateTaskRequestImpl(
      {this.name,
      this.projectId,
      this.description,
      this.autoCalculated,
      this.type,
      this.status,
      this.category,
      this.deadline});

  factory _$UpdateTaskRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$UpdateTaskRequestImplFromJson(json);

  @override
  final String? name;
  @override
  final int? projectId;
  @override
  final String? description;
  @override
  final bool? autoCalculated;
  @override
  final int? type;
  @override
  final int? status;
  @override
  final String? category;
  @override
  final String? deadline;

  @override
  String toString() {
    return 'UpdateTaskRequest(name: $name, projectId: $projectId, description: $description, autoCalculated: $autoCalculated, type: $type, status: $status, category: $category, deadline: $deadline)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UpdateTaskRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.projectId, projectId) ||
                other.projectId == projectId) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.autoCalculated, autoCalculated) ||
                other.autoCalculated == autoCalculated) &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.category, category) ||
                other.category == category) &&
            (identical(other.deadline, deadline) ||
                other.deadline == deadline));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, name, projectId, description,
      autoCalculated, type, status, category, deadline);

  /// Create a copy of UpdateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UpdateTaskRequestImplCopyWith<_$UpdateTaskRequestImpl> get copyWith =>
      __$$UpdateTaskRequestImplCopyWithImpl<_$UpdateTaskRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UpdateTaskRequestImplToJson(
      this,
    );
  }
}

abstract class _UpdateTaskRequest implements UpdateTaskRequest {
  const factory _UpdateTaskRequest(
      {final String? name,
      final int? projectId,
      final String? description,
      final bool? autoCalculated,
      final int? type,
      final int? status,
      final String? category,
      final String? deadline}) = _$UpdateTaskRequestImpl;

  factory _UpdateTaskRequest.fromJson(Map<String, dynamic> json) =
      _$UpdateTaskRequestImpl.fromJson;

  @override
  String? get name;
  @override
  int? get projectId;
  @override
  String? get description;
  @override
  bool? get autoCalculated;
  @override
  int? get type;
  @override
  int? get status;
  @override
  String? get category;
  @override
  String? get deadline;

  /// Create a copy of UpdateTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UpdateTaskRequestImplCopyWith<_$UpdateTaskRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

FinishTaskRequest _$FinishTaskRequestFromJson(Map<String, dynamic> json) {
  return _FinishTaskRequest.fromJson(json);
}

/// @nodoc
mixin _$FinishTaskRequest {
  String get time => throw _privateConstructorUsedError;

  /// Serializes this FinishTaskRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of FinishTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $FinishTaskRequestCopyWith<FinishTaskRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $FinishTaskRequestCopyWith<$Res> {
  factory $FinishTaskRequestCopyWith(
          FinishTaskRequest value, $Res Function(FinishTaskRequest) then) =
      _$FinishTaskRequestCopyWithImpl<$Res, FinishTaskRequest>;
  @useResult
  $Res call({String time});
}

/// @nodoc
class _$FinishTaskRequestCopyWithImpl<$Res, $Val extends FinishTaskRequest>
    implements $FinishTaskRequestCopyWith<$Res> {
  _$FinishTaskRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of FinishTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? time = null,
  }) {
    return _then(_value.copyWith(
      time: null == time
          ? _value.time
          : time // ignore: cast_nullable_to_non_nullable
              as String,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$FinishTaskRequestImplCopyWith<$Res>
    implements $FinishTaskRequestCopyWith<$Res> {
  factory _$$FinishTaskRequestImplCopyWith(_$FinishTaskRequestImpl value,
          $Res Function(_$FinishTaskRequestImpl) then) =
      __$$FinishTaskRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String time});
}

/// @nodoc
class __$$FinishTaskRequestImplCopyWithImpl<$Res>
    extends _$FinishTaskRequestCopyWithImpl<$Res, _$FinishTaskRequestImpl>
    implements _$$FinishTaskRequestImplCopyWith<$Res> {
  __$$FinishTaskRequestImplCopyWithImpl(_$FinishTaskRequestImpl _value,
      $Res Function(_$FinishTaskRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of FinishTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? time = null,
  }) {
    return _then(_$FinishTaskRequestImpl(
      time: null == time
          ? _value.time
          : time // ignore: cast_nullable_to_non_nullable
              as String,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$FinishTaskRequestImpl implements _FinishTaskRequest {
  const _$FinishTaskRequestImpl({required this.time});

  factory _$FinishTaskRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$FinishTaskRequestImplFromJson(json);

  @override
  final String time;

  @override
  String toString() {
    return 'FinishTaskRequest(time: $time)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$FinishTaskRequestImpl &&
            (identical(other.time, time) || other.time == time));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, time);

  /// Create a copy of FinishTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$FinishTaskRequestImplCopyWith<_$FinishTaskRequestImpl> get copyWith =>
      __$$FinishTaskRequestImplCopyWithImpl<_$FinishTaskRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$FinishTaskRequestImplToJson(
      this,
    );
  }
}

abstract class _FinishTaskRequest implements FinishTaskRequest {
  const factory _FinishTaskRequest({required final String time}) =
      _$FinishTaskRequestImpl;

  factory _FinishTaskRequest.fromJson(Map<String, dynamic> json) =
      _$FinishTaskRequestImpl.fromJson;

  @override
  String get time;

  /// Create a copy of FinishTaskRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$FinishTaskRequestImplCopyWith<_$FinishTaskRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateTaskPaymentRequest _$CreateTaskPaymentRequestFromJson(
    Map<String, dynamic> json) {
  return _CreateTaskPaymentRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateTaskPaymentRequest {
  int get budgetId => throw _privateConstructorUsedError;
  double get amount => throw _privateConstructorUsedError;

  /// Serializes this CreateTaskPaymentRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of CreateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CreateTaskPaymentRequestCopyWith<CreateTaskPaymentRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateTaskPaymentRequestCopyWith<$Res> {
  factory $CreateTaskPaymentRequestCopyWith(CreateTaskPaymentRequest value,
          $Res Function(CreateTaskPaymentRequest) then) =
      _$CreateTaskPaymentRequestCopyWithImpl<$Res, CreateTaskPaymentRequest>;
  @useResult
  $Res call({int budgetId, double amount});
}

/// @nodoc
class _$CreateTaskPaymentRequestCopyWithImpl<$Res,
        $Val extends CreateTaskPaymentRequest>
    implements $CreateTaskPaymentRequestCopyWith<$Res> {
  _$CreateTaskPaymentRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CreateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? budgetId = null,
    Object? amount = null,
  }) {
    return _then(_value.copyWith(
      budgetId: null == budgetId
          ? _value.budgetId
          : budgetId // ignore: cast_nullable_to_non_nullable
              as int,
      amount: null == amount
          ? _value.amount
          : amount // ignore: cast_nullable_to_non_nullable
              as double,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CreateTaskPaymentRequestImplCopyWith<$Res>
    implements $CreateTaskPaymentRequestCopyWith<$Res> {
  factory _$$CreateTaskPaymentRequestImplCopyWith(
          _$CreateTaskPaymentRequestImpl value,
          $Res Function(_$CreateTaskPaymentRequestImpl) then) =
      __$$CreateTaskPaymentRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({int budgetId, double amount});
}

/// @nodoc
class __$$CreateTaskPaymentRequestImplCopyWithImpl<$Res>
    extends _$CreateTaskPaymentRequestCopyWithImpl<$Res,
        _$CreateTaskPaymentRequestImpl>
    implements _$$CreateTaskPaymentRequestImplCopyWith<$Res> {
  __$$CreateTaskPaymentRequestImplCopyWithImpl(
      _$CreateTaskPaymentRequestImpl _value,
      $Res Function(_$CreateTaskPaymentRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of CreateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? budgetId = null,
    Object? amount = null,
  }) {
    return _then(_$CreateTaskPaymentRequestImpl(
      budgetId: null == budgetId
          ? _value.budgetId
          : budgetId // ignore: cast_nullable_to_non_nullable
              as int,
      amount: null == amount
          ? _value.amount
          : amount // ignore: cast_nullable_to_non_nullable
              as double,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CreateTaskPaymentRequestImpl implements _CreateTaskPaymentRequest {
  const _$CreateTaskPaymentRequestImpl(
      {required this.budgetId, required this.amount});

  factory _$CreateTaskPaymentRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$CreateTaskPaymentRequestImplFromJson(json);

  @override
  final int budgetId;
  @override
  final double amount;

  @override
  String toString() {
    return 'CreateTaskPaymentRequest(budgetId: $budgetId, amount: $amount)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateTaskPaymentRequestImpl &&
            (identical(other.budgetId, budgetId) ||
                other.budgetId == budgetId) &&
            (identical(other.amount, amount) || other.amount == amount));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, budgetId, amount);

  /// Create a copy of CreateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateTaskPaymentRequestImplCopyWith<_$CreateTaskPaymentRequestImpl>
      get copyWith => __$$CreateTaskPaymentRequestImplCopyWithImpl<
          _$CreateTaskPaymentRequestImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateTaskPaymentRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateTaskPaymentRequest implements CreateTaskPaymentRequest {
  const factory _CreateTaskPaymentRequest(
      {required final int budgetId,
      required final double amount}) = _$CreateTaskPaymentRequestImpl;

  factory _CreateTaskPaymentRequest.fromJson(Map<String, dynamic> json) =
      _$CreateTaskPaymentRequestImpl.fromJson;

  @override
  int get budgetId;
  @override
  double get amount;

  /// Create a copy of CreateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CreateTaskPaymentRequestImplCopyWith<_$CreateTaskPaymentRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

UpdateTaskPaymentRequest _$UpdateTaskPaymentRequestFromJson(
    Map<String, dynamic> json) {
  return _UpdateTaskPaymentRequest.fromJson(json);
}

/// @nodoc
mixin _$UpdateTaskPaymentRequest {
  double get amount => throw _privateConstructorUsedError;

  /// Serializes this UpdateTaskPaymentRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of UpdateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UpdateTaskPaymentRequestCopyWith<UpdateTaskPaymentRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UpdateTaskPaymentRequestCopyWith<$Res> {
  factory $UpdateTaskPaymentRequestCopyWith(UpdateTaskPaymentRequest value,
          $Res Function(UpdateTaskPaymentRequest) then) =
      _$UpdateTaskPaymentRequestCopyWithImpl<$Res, UpdateTaskPaymentRequest>;
  @useResult
  $Res call({double amount});
}

/// @nodoc
class _$UpdateTaskPaymentRequestCopyWithImpl<$Res,
        $Val extends UpdateTaskPaymentRequest>
    implements $UpdateTaskPaymentRequestCopyWith<$Res> {
  _$UpdateTaskPaymentRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of UpdateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? amount = null,
  }) {
    return _then(_value.copyWith(
      amount: null == amount
          ? _value.amount
          : amount // ignore: cast_nullable_to_non_nullable
              as double,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$UpdateTaskPaymentRequestImplCopyWith<$Res>
    implements $UpdateTaskPaymentRequestCopyWith<$Res> {
  factory _$$UpdateTaskPaymentRequestImplCopyWith(
          _$UpdateTaskPaymentRequestImpl value,
          $Res Function(_$UpdateTaskPaymentRequestImpl) then) =
      __$$UpdateTaskPaymentRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({double amount});
}

/// @nodoc
class __$$UpdateTaskPaymentRequestImplCopyWithImpl<$Res>
    extends _$UpdateTaskPaymentRequestCopyWithImpl<$Res,
        _$UpdateTaskPaymentRequestImpl>
    implements _$$UpdateTaskPaymentRequestImplCopyWith<$Res> {
  __$$UpdateTaskPaymentRequestImplCopyWithImpl(
      _$UpdateTaskPaymentRequestImpl _value,
      $Res Function(_$UpdateTaskPaymentRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of UpdateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? amount = null,
  }) {
    return _then(_$UpdateTaskPaymentRequestImpl(
      amount: null == amount
          ? _value.amount
          : amount // ignore: cast_nullable_to_non_nullable
              as double,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UpdateTaskPaymentRequestImpl implements _UpdateTaskPaymentRequest {
  const _$UpdateTaskPaymentRequestImpl({required this.amount});

  factory _$UpdateTaskPaymentRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$UpdateTaskPaymentRequestImplFromJson(json);

  @override
  final double amount;

  @override
  String toString() {
    return 'UpdateTaskPaymentRequest(amount: $amount)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UpdateTaskPaymentRequestImpl &&
            (identical(other.amount, amount) || other.amount == amount));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, amount);

  /// Create a copy of UpdateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UpdateTaskPaymentRequestImplCopyWith<_$UpdateTaskPaymentRequestImpl>
      get copyWith => __$$UpdateTaskPaymentRequestImplCopyWithImpl<
          _$UpdateTaskPaymentRequestImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UpdateTaskPaymentRequestImplToJson(
      this,
    );
  }
}

abstract class _UpdateTaskPaymentRequest implements UpdateTaskPaymentRequest {
  const factory _UpdateTaskPaymentRequest({required final double amount}) =
      _$UpdateTaskPaymentRequestImpl;

  factory _UpdateTaskPaymentRequest.fromJson(Map<String, dynamic> json) =
      _$UpdateTaskPaymentRequestImpl.fromJson;

  @override
  double get amount;

  /// Create a copy of UpdateTaskPaymentRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UpdateTaskPaymentRequestImplCopyWith<_$UpdateTaskPaymentRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

PrerequisitesRequest _$PrerequisitesRequestFromJson(Map<String, dynamic> json) {
  return _PrerequisitesRequest.fromJson(json);
}

/// @nodoc
mixin _$PrerequisitesRequest {
  int get prerequisiteId => throw _privateConstructorUsedError;

  /// Serializes this PrerequisitesRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of PrerequisitesRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $PrerequisitesRequestCopyWith<PrerequisitesRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $PrerequisitesRequestCopyWith<$Res> {
  factory $PrerequisitesRequestCopyWith(PrerequisitesRequest value,
          $Res Function(PrerequisitesRequest) then) =
      _$PrerequisitesRequestCopyWithImpl<$Res, PrerequisitesRequest>;
  @useResult
  $Res call({int prerequisiteId});
}

/// @nodoc
class _$PrerequisitesRequestCopyWithImpl<$Res,
        $Val extends PrerequisitesRequest>
    implements $PrerequisitesRequestCopyWith<$Res> {
  _$PrerequisitesRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of PrerequisitesRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? prerequisiteId = null,
  }) {
    return _then(_value.copyWith(
      prerequisiteId: null == prerequisiteId
          ? _value.prerequisiteId
          : prerequisiteId // ignore: cast_nullable_to_non_nullable
              as int,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$PrerequisitesRequestImplCopyWith<$Res>
    implements $PrerequisitesRequestCopyWith<$Res> {
  factory _$$PrerequisitesRequestImplCopyWith(_$PrerequisitesRequestImpl value,
          $Res Function(_$PrerequisitesRequestImpl) then) =
      __$$PrerequisitesRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({int prerequisiteId});
}

/// @nodoc
class __$$PrerequisitesRequestImplCopyWithImpl<$Res>
    extends _$PrerequisitesRequestCopyWithImpl<$Res, _$PrerequisitesRequestImpl>
    implements _$$PrerequisitesRequestImplCopyWith<$Res> {
  __$$PrerequisitesRequestImplCopyWithImpl(_$PrerequisitesRequestImpl _value,
      $Res Function(_$PrerequisitesRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of PrerequisitesRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? prerequisiteId = null,
  }) {
    return _then(_$PrerequisitesRequestImpl(
      prerequisiteId: null == prerequisiteId
          ? _value.prerequisiteId
          : prerequisiteId // ignore: cast_nullable_to_non_nullable
              as int,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$PrerequisitesRequestImpl implements _PrerequisitesRequest {
  const _$PrerequisitesRequestImpl({required this.prerequisiteId});

  factory _$PrerequisitesRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$PrerequisitesRequestImplFromJson(json);

  @override
  final int prerequisiteId;

  @override
  String toString() {
    return 'PrerequisitesRequest(prerequisiteId: $prerequisiteId)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$PrerequisitesRequestImpl &&
            (identical(other.prerequisiteId, prerequisiteId) ||
                other.prerequisiteId == prerequisiteId));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, prerequisiteId);

  /// Create a copy of PrerequisitesRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$PrerequisitesRequestImplCopyWith<_$PrerequisitesRequestImpl>
      get copyWith =>
          __$$PrerequisitesRequestImplCopyWithImpl<_$PrerequisitesRequestImpl>(
              this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$PrerequisitesRequestImplToJson(
      this,
    );
  }
}

abstract class _PrerequisitesRequest implements PrerequisitesRequest {
  const factory _PrerequisitesRequest({required final int prerequisiteId}) =
      _$PrerequisitesRequestImpl;

  factory _PrerequisitesRequest.fromJson(Map<String, dynamic> json) =
      _$PrerequisitesRequestImpl.fromJson;

  @override
  int get prerequisiteId;

  /// Create a copy of PrerequisitesRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$PrerequisitesRequestImplCopyWith<_$PrerequisitesRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

Dependency _$DependencyFromJson(Map<String, dynamic> json) {
  return _Dependency.fromJson(json);
}

/// @nodoc
mixin _$Dependency {
  int get prerequisiteId => throw _privateConstructorUsedError;
  int get taskId => throw _privateConstructorUsedError;

  /// Serializes this Dependency to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Dependency
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $DependencyCopyWith<Dependency> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $DependencyCopyWith<$Res> {
  factory $DependencyCopyWith(
          Dependency value, $Res Function(Dependency) then) =
      _$DependencyCopyWithImpl<$Res, Dependency>;
  @useResult
  $Res call({int prerequisiteId, int taskId});
}

/// @nodoc
class _$DependencyCopyWithImpl<$Res, $Val extends Dependency>
    implements $DependencyCopyWith<$Res> {
  _$DependencyCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Dependency
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? prerequisiteId = null,
    Object? taskId = null,
  }) {
    return _then(_value.copyWith(
      prerequisiteId: null == prerequisiteId
          ? _value.prerequisiteId
          : prerequisiteId // ignore: cast_nullable_to_non_nullable
              as int,
      taskId: null == taskId
          ? _value.taskId
          : taskId // ignore: cast_nullable_to_non_nullable
              as int,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$DependencyImplCopyWith<$Res>
    implements $DependencyCopyWith<$Res> {
  factory _$$DependencyImplCopyWith(
          _$DependencyImpl value, $Res Function(_$DependencyImpl) then) =
      __$$DependencyImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({int prerequisiteId, int taskId});
}

/// @nodoc
class __$$DependencyImplCopyWithImpl<$Res>
    extends _$DependencyCopyWithImpl<$Res, _$DependencyImpl>
    implements _$$DependencyImplCopyWith<$Res> {
  __$$DependencyImplCopyWithImpl(
      _$DependencyImpl _value, $Res Function(_$DependencyImpl) _then)
      : super(_value, _then);

  /// Create a copy of Dependency
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? prerequisiteId = null,
    Object? taskId = null,
  }) {
    return _then(_$DependencyImpl(
      prerequisiteId: null == prerequisiteId
          ? _value.prerequisiteId
          : prerequisiteId // ignore: cast_nullable_to_non_nullable
              as int,
      taskId: null == taskId
          ? _value.taskId
          : taskId // ignore: cast_nullable_to_non_nullable
              as int,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$DependencyImpl implements _Dependency {
  const _$DependencyImpl({required this.prerequisiteId, required this.taskId});

  factory _$DependencyImpl.fromJson(Map<String, dynamic> json) =>
      _$$DependencyImplFromJson(json);

  @override
  final int prerequisiteId;
  @override
  final int taskId;

  @override
  String toString() {
    return 'Dependency(prerequisiteId: $prerequisiteId, taskId: $taskId)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$DependencyImpl &&
            (identical(other.prerequisiteId, prerequisiteId) ||
                other.prerequisiteId == prerequisiteId) &&
            (identical(other.taskId, taskId) || other.taskId == taskId));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, prerequisiteId, taskId);

  /// Create a copy of Dependency
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$DependencyImplCopyWith<_$DependencyImpl> get copyWith =>
      __$$DependencyImplCopyWithImpl<_$DependencyImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$DependencyImplToJson(
      this,
    );
  }
}

abstract class _Dependency implements Dependency {
  const factory _Dependency(
      {required final int prerequisiteId,
      required final int taskId}) = _$DependencyImpl;

  factory _Dependency.fromJson(Map<String, dynamic> json) =
      _$DependencyImpl.fromJson;

  @override
  int get prerequisiteId;
  @override
  int get taskId;

  /// Create a copy of Dependency
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$DependencyImplCopyWith<_$DependencyImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

TaskListResponse _$TaskListResponseFromJson(Map<String, dynamic> json) {
  return _TaskListResponse.fromJson(json);
}

/// @nodoc
mixin _$TaskListResponse {
  int get page => throw _privateConstructorUsedError;
  int get total => throw _privateConstructorUsedError;
  int get size => throw _privateConstructorUsedError;
  List<Task> get list => throw _privateConstructorUsedError;

  /// Serializes this TaskListResponse to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of TaskListResponse
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $TaskListResponseCopyWith<TaskListResponse> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $TaskListResponseCopyWith<$Res> {
  factory $TaskListResponseCopyWith(
          TaskListResponse value, $Res Function(TaskListResponse) then) =
      _$TaskListResponseCopyWithImpl<$Res, TaskListResponse>;
  @useResult
  $Res call({int page, int total, int size, List<Task> list});
}

/// @nodoc
class _$TaskListResponseCopyWithImpl<$Res, $Val extends TaskListResponse>
    implements $TaskListResponseCopyWith<$Res> {
  _$TaskListResponseCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of TaskListResponse
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? page = null,
    Object? total = null,
    Object? size = null,
    Object? list = null,
  }) {
    return _then(_value.copyWith(
      page: null == page
          ? _value.page
          : page // ignore: cast_nullable_to_non_nullable
              as int,
      total: null == total
          ? _value.total
          : total // ignore: cast_nullable_to_non_nullable
              as int,
      size: null == size
          ? _value.size
          : size // ignore: cast_nullable_to_non_nullable
              as int,
      list: null == list
          ? _value.list
          : list // ignore: cast_nullable_to_non_nullable
              as List<Task>,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$TaskListResponseImplCopyWith<$Res>
    implements $TaskListResponseCopyWith<$Res> {
  factory _$$TaskListResponseImplCopyWith(_$TaskListResponseImpl value,
          $Res Function(_$TaskListResponseImpl) then) =
      __$$TaskListResponseImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({int page, int total, int size, List<Task> list});
}

/// @nodoc
class __$$TaskListResponseImplCopyWithImpl<$Res>
    extends _$TaskListResponseCopyWithImpl<$Res, _$TaskListResponseImpl>
    implements _$$TaskListResponseImplCopyWith<$Res> {
  __$$TaskListResponseImplCopyWithImpl(_$TaskListResponseImpl _value,
      $Res Function(_$TaskListResponseImpl) _then)
      : super(_value, _then);

  /// Create a copy of TaskListResponse
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? page = null,
    Object? total = null,
    Object? size = null,
    Object? list = null,
  }) {
    return _then(_$TaskListResponseImpl(
      page: null == page
          ? _value.page
          : page // ignore: cast_nullable_to_non_nullable
              as int,
      total: null == total
          ? _value.total
          : total // ignore: cast_nullable_to_non_nullable
              as int,
      size: null == size
          ? _value.size
          : size // ignore: cast_nullable_to_non_nullable
              as int,
      list: null == list
          ? _value._list
          : list // ignore: cast_nullable_to_non_nullable
              as List<Task>,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$TaskListResponseImpl implements _TaskListResponse {
  const _$TaskListResponseImpl(
      {required this.page,
      required this.total,
      required this.size,
      required final List<Task> list})
      : _list = list;

  factory _$TaskListResponseImpl.fromJson(Map<String, dynamic> json) =>
      _$$TaskListResponseImplFromJson(json);

  @override
  final int page;
  @override
  final int total;
  @override
  final int size;
  final List<Task> _list;
  @override
  List<Task> get list {
    if (_list is EqualUnmodifiableListView) return _list;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_list);
  }

  @override
  String toString() {
    return 'TaskListResponse(page: $page, total: $total, size: $size, list: $list)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$TaskListResponseImpl &&
            (identical(other.page, page) || other.page == page) &&
            (identical(other.total, total) || other.total == total) &&
            (identical(other.size, size) || other.size == size) &&
            const DeepCollectionEquality().equals(other._list, _list));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, page, total, size,
      const DeepCollectionEquality().hash(_list));

  /// Create a copy of TaskListResponse
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$TaskListResponseImplCopyWith<_$TaskListResponseImpl> get copyWith =>
      __$$TaskListResponseImplCopyWithImpl<_$TaskListResponseImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$TaskListResponseImplToJson(
      this,
    );
  }
}

abstract class _TaskListResponse implements TaskListResponse {
  const factory _TaskListResponse(
      {required final int page,
      required final int total,
      required final int size,
      required final List<Task> list}) = _$TaskListResponseImpl;

  factory _TaskListResponse.fromJson(Map<String, dynamic> json) =
      _$TaskListResponseImpl.fromJson;

  @override
  int get page;
  @override
  int get total;
  @override
  int get size;
  @override
  List<Task> get list;

  /// Create a copy of TaskListResponse
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$TaskListResponseImplCopyWith<_$TaskListResponseImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
