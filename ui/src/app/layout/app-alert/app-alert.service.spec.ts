import { TestBed } from '@angular/core/testing';

import { AppAlertService } from './app-alert.service';

describe('AppAlertService', () => {
  let service: AppAlertService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AppAlertService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
