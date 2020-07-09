import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneDetailComponent } from './zone-detail.component';

describe('ZoneDetailComponent', () => {
  let component: ZoneDetailComponent;
  let fixture: ComponentFixture<ZoneDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ZoneDetailComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ZoneDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
