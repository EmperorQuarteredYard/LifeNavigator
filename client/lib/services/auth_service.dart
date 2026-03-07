import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/models.dart';
import 'api_client.dart';

class AuthService {
  final ApiClient _client;

  AuthService(this._client);

  Future<AuthResponse> login(String username, String password) async {
    final response = await _client.dio.post(
      '/auth/login',
      data: LoginRequest(username: username, password: password).toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    final authResponse = AuthResponse.fromJson(apiResponse.data!);
    _client.setTokens(authResponse.accessToken, authResponse.refreshToken);
    return authResponse;
  }

  Future<AuthResponse> register({
    required String username,
    required String password,
    String? nickname,
    String? email,
    String? phone,
    required String inviteCode,
  }) async {
    final response = await _client.dio.post(
      '/auth/register',
      data: RegisterRequest(
        username: username,
        password: password,
        nickname: nickname,
        email: email,
        phone: phone,
        inviteCode: inviteCode,
      ).toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    final authResponse = AuthResponse.fromJson(apiResponse.data!);
    _client.setTokens(authResponse.accessToken, authResponse.refreshToken);
    return authResponse;
  }

  Future<RefreshResponse> refreshToken(String refreshToken) async {
    final response = await _client.dio.post(
      '/auth/refresh',
      data: RefreshRequest(refreshToken: refreshToken).toJson(),
    );
    final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
      response.data,
      (json) => json as Map<String, dynamic>,
    );
    if (apiResponse.code != 0) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    final refreshResponse = RefreshResponse.fromJson(apiResponse.data!);
    _client.setTokens(refreshResponse.accessToken, refreshResponse.refreshToken);
    return refreshResponse;
  }

  void logout() {
    _client.clearTokens();
  }
}

final authServiceProvider = Provider<AuthService>((ref) {
  return AuthService(ref.read(apiClientProvider));
});
