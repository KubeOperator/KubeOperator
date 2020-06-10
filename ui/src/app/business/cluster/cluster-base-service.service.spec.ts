import { TestBed } from '@angular/core/testing';

import { ClusterBaseServiceService } from './cluster-base-service.service';

describe('ClusterBaseServiceService', () => {
  let service: ClusterBaseServiceService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ClusterBaseServiceService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
