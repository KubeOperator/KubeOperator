import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RegionListComponent } from './region-list.component';

describe('RegionListComponent', () => {
  let component: RegionListComponent;
  let fixture: ComponentFixture<RegionListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RegionListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RegionListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
