import { TestBed } from '@angular/core/testing';

import { VmConfigService } from './vm-config.service';

describe('VmConfigService', () => {
  let service: VmConfigService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(VmConfigService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
