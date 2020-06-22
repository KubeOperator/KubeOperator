import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MonitorEnableComponent } from './monitor-enable.component';

describe('MonitorEnableComponent', () => {
  let component: MonitorEnableComponent;
  let fixture: ComponentFixture<MonitorEnableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MonitorEnableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MonitorEnableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
