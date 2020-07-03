import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PrometheusEnableComponent } from './prometheus-enable.component';

describe('PrometheusEnableComponent', () => {
  let component: PrometheusEnableComponent;
  let fixture: ComponentFixture<PrometheusEnableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PrometheusEnableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PrometheusEnableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
