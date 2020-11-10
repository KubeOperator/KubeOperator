import { TestBed } from '@angular/core/testing';

import { BusinessResolverService } from './business-resolver.service';

describe('BusinessResolverService', () => {
  let service: BusinessResolverService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BusinessResolverService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
