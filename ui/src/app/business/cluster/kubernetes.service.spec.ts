import { TestBed } from '@angular/core/testing';

import { KubernetesService } from './kubernetes.service';

describe('KubernetesService', () => {
  let service: KubernetesService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(KubernetesService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
