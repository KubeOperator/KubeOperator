import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IngressInstallComponent } from './ingress-install.component';

describe('IngressInstallComponent', () => {
  let component: IngressInstallComponent;
  let fixture: ComponentFixture<IngressInstallComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IngressInstallComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IngressInstallComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
