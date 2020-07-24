import { TestBed } from '@angular/core/testing';

import { ProjectMemberService } from './project-member.service';

describe('ProjectMemberService', () => {
  let service: ProjectMemberService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProjectMemberService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
