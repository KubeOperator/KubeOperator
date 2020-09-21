import { TestBed } from '@angular/core/testing';

import { UserSubscribeService } from './user-subscribe.service';

describe('UserSubscribeService', () => {
  let service: UserSubscribeService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(UserSubscribeService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
