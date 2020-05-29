import { TestBed } from '@angular/core/testing';

import { CommonAlertService } from './common-alert.service';

describe('CommonAlertService', () => {
  let service: CommonAlertService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CommonAlertService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
