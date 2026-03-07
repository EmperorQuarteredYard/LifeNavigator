import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/models.dart';
import 'api_client.dart';

class AccountService {
  final ApiClient _client;

  AccountService(this._client);

  Future<Account> createAccount(CreateAccountRequest request) async {
    final response = await _client.dio.post(
      '/accounts',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Account.fromJson(apiResponse.data!);
  }

  Future<List<Account>> getAccounts() async {
    final response = await _client.dio.get('/accounts');
    final apiResponse = ApiResponse<List<dynamic>>.fromJson(
      response.data,
      (json) => json as List<dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!
        .map((e) => Account.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Account> getAccount(int id) async {
    final response = await _client.dio.get('/accounts/$id');
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Account.fromJson(apiResponse.data!);
  }

  Future<void> deleteAccount(int id) async {
    final response = await _client.dio.delete('/accounts/$id');
    final apiResponse = ApiResponse<dynamic>.fromJson(
      response.data,
      (json) => json,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }

  Future<Account> adjustBalance(int id, double amount) async {
    final response = await _client.dio.post(
      '/accounts/$id/adjust',
      data: AdjustBalanceRequest(amount: amount).toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return Account.fromJson(apiResponse.data!);
  }
}

final accountServiceProvider = Provider<AccountService>((ref) {
  return AccountService(ref.read(apiClientProvider));
});
