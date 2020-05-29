import { TestBed } from '@angular/core/testing';

import { ModalAlertService } from './modal-alert.service';

describe('ModalAlertService', () => {
  let service: ModalAlertService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModalAlertService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
