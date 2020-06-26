import { TestBed } from '@angular/core/testing';

import { CloudProviderService } from './cloud-provider.service';

describe('CloudProviderService', () => {
  let service: CloudProviderService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CloudProviderService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
