import { TestBed } from '@angular/core/testing';

import { IpService } from './ip.service';

describe('IpService', () => {
  let service: IpService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(IpService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
