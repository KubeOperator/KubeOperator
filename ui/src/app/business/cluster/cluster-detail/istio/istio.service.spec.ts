import { TestBed } from '@angular/core/testing';

import { IstioService } from './istio.service';

describe('IstioService', () => {
  let service: IstioService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(IstioService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
