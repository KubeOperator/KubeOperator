import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardEnableComponent } from './dashboard-enable.component';

describe('DashboardEnableComponent', () => {
  let component: DashboardEnableComponent;
  let fixture: ComponentFixture<DashboardEnableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DashboardEnableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DashboardEnableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
