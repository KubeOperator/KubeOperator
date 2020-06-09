import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LoggingQueryComponent } from './logging-query.component';

describe('LoggingQueryComponent', () => {
  let component: LoggingQueryComponent;
  let fixture: ComponentFixture<LoggingQueryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LoggingQueryComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoggingQueryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
