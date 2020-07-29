import { TestBed } from '@angular/core/testing';

import { ProjectClusterResolverService } from './project-cluster-resolver.service';

describe('ProjectClusterResolverService', () => {
  let service: ProjectClusterResolverService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProjectClusterResolverService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
