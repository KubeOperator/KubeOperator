import { TestBed } from '@angular/core/testing';

import { ProjectResourceService } from './project-resource.service';

describe('ProjectResourceService', () => {
  let service: ProjectResourceService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProjectResourceService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
