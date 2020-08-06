import { BackupAccountStatusPipe } from './backup-account-status.pipe';

describe('BackupAccountStatusPipe', () => {
  it('create an instance', () => {
    const pipe = new BackupAccountStatusPipe();
    expect(pipe).toBeTruthy();
  });
});
