import { TestBed } from '@angular/core/testing';

import { StorageProvisionerService } from './storage-provisioner.service';

describe('StorageProvisionerService', () => {
  let service: StorageProvisionerService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(StorageProvisionerService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
