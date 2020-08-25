import { TestBed } from '@angular/core/testing';

import { ManifestService } from './manifest.service';

describe('ManifestService', () => {
  let service: ManifestService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ManifestService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
