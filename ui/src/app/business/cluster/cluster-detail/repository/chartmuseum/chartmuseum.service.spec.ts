import { TestBed } from '@angular/core/testing';

import { ChartmuseumService } from './chartmuseum.service';

describe('ChartmuseumService', () => {
  let service: ChartmuseumService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ChartmuseumService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
