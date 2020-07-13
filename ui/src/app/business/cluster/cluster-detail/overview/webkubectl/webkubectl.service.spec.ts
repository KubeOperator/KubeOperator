import { TestBed } from '@angular/core/testing';

import { WebkubectlService } from './webkubectl.service';

describe('WebkubectlService', () => {
  let service: WebkubectlService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(WebkubectlService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
