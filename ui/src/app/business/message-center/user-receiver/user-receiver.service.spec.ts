import { TestBed } from '@angular/core/testing';

import { UserReceiverService } from './user-receiver.service';

describe('UserReceiverService', () => {
  let service: UserReceiverService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(UserReceiverService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
