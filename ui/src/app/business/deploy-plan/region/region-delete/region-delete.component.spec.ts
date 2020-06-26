import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RegionDeleteComponent } from './region-delete.component';

describe('RegionDeleteComponent', () => {
  let component: RegionDeleteComponent;
  let fixture: ComponentFixture<RegionDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RegionDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RegionDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
