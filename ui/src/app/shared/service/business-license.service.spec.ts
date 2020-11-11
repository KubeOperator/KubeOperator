import { TestBed } from '@angular/core/testing';

import { BusinessLicenseService } from './business-license.service';

describe('BusinessLicenseService', () => {
  let service: BusinessLicenseService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BusinessLicenseService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
