import { TestBed } from '@angular/core/testing';

import { ClusterLoggerService } from './cluster-logger.service';

describe('ClusterLoggerService', () => {
  let service: ClusterLoggerService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ClusterLoggerService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
