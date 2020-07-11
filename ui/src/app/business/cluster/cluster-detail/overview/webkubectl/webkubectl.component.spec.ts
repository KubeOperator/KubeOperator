import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WebkubectlComponent } from './webkubectl.component';

describe('WebkubectlComponent', () => {
  let component: WebkubectlComponent;
  let fixture: ComponentFixture<WebkubectlComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WebkubectlComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WebkubectlComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
