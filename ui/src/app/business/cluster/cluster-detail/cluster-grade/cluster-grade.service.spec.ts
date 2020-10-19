import { TestBed } from '@angular/core/testing';

import { ClusterGradeService } from './cluster-grade.service';

describe('ClusterGradeService', () => {
  let service: ClusterGradeService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ClusterGradeService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
