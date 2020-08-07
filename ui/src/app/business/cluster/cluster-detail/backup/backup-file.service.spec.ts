import { TestBed } from '@angular/core/testing';

import { BackupFileService } from './backup-file.service';

describe('BackupFileService', () => {
  let service: BackupFileService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BackupFileService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
