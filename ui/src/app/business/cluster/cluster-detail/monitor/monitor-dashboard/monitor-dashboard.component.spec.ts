import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MonitorDashboardComponent } from './monitor-dashboard.component';

describe('MonitorDashboardComponent', () => {
  let component: MonitorDashboardComponent;
  let fixture: ComponentFixture<MonitorDashboardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MonitorDashboardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MonitorDashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
