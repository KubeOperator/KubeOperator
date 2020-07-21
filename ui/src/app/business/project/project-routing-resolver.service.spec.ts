import { TestBed } from '@angular/core/testing';

import { ProjectRoutingResolverService } from './project-routing-resolver.service';

describe('ProjectRoutingResolverService', () => {
  let service: ProjectRoutingResolverService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProjectRoutingResolverService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
