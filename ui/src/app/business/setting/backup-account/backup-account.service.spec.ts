import { TestBed } from '@angular/core/testing';

import { BackupAccountService } from './backup-account.service';

describe('BackupAccountService', () => {
  let service: BackupAccountService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BackupAccountService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
