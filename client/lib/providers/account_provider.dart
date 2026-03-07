import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/models.dart';
import '../services/services.dart';

class AccountState {
  final List<Account> accounts;
  final Account? selectedAccount;
  final bool isLoading;
  final String? error;

  const AccountState({
    this.accounts = const [],
    this.selectedAccount,
    this.isLoading = false,
    this.error,
  });

  AccountState copyWith({
    List<Account>? accounts,
    Account? selectedAccount,
    bool? isLoading,
    String? error,
  }) {
    return AccountState(
      accounts: accounts ?? this.accounts,
      selectedAccount: selectedAccount ?? this.selectedAccount,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }

  double get totalBalance =>
      accounts.fold(0, (sum, account) => sum + account.balance);
}

class AccountNotifier extends StateNotifier<AccountState> {
  final AccountService _accountService;

  AccountNotifier(this._accountService) : super(const AccountState());

  Future<void> loadAccounts() async {
    if (state.isLoading) return;

    state = state.copyWith(isLoading: true, error: null);

    try {
      final accounts = await _accountService.getAccounts();
      state = state.copyWith(
        accounts: accounts,
        isLoading: false,
      );
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<bool> createAccount(CreateAccountRequest request) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final account = await _accountService.createAccount(request);
      state = state.copyWith(
        accounts: [...state.accounts, account],
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

  Future<bool> deleteAccount(int id) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _accountService.deleteAccount(id);
      final newAccounts = state.accounts.where((a) => a.id != id).toList();
      state = state.copyWith(
        accounts: newAccounts,
        selectedAccount:
            state.selectedAccount?.id == id ? null : state.selectedAccount,
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

  Future<bool> adjustBalance(int id, double amount) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final updatedAccount = await _accountService.adjustBalance(id, amount);
      final newAccounts = state.accounts.map((a) {
        return a.id == id ? updatedAccount : a;
      }).toList();
      state = state.copyWith(
        accounts: newAccounts,
        selectedAccount: state.selectedAccount?.id == id
            ? updatedAccount
            : state.selectedAccount,
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

  void selectAccount(int id) {
    final account = state.accounts.where((a) => a.id == id).firstOrNull;
    state = state.copyWith(selectedAccount: account);
  }

  void clearSelection() {
    state = state.copyWith(selectedAccount: null);
  }
}

final accountStateProvider =
    StateNotifierProvider<AccountNotifier, AccountState>((ref) {
  return AccountNotifier(ref.read(accountServiceProvider));
});

final accountByIdProvider = Provider.family<Account?, int>((ref, id) {
  final state = ref.watch(accountStateProvider);
  return state.accounts.where((a) => a.id == id).firstOrNull;
});
