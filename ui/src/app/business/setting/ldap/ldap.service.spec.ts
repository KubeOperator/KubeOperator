import { TestBed } from '@angular/core/testing';

import { LdapService } from './ldap.service';

describe('LdapService', () => {
  let service: LdapService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(LdapService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
