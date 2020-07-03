import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LoggingEnableComponent } from './logging-enable.component';

describe('LoggingEnableComponent', () => {
  let component: LoggingEnableComponent;
  let fixture: ComponentFixture<LoggingEnableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LoggingEnableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoggingEnableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
